package server

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	arkv1 "github.com/fedragon/ark/gen/ark/v1"
	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/metrics"

	connect_go "github.com/bufbuild/connect-go"
)

type Handler struct {
	Repo        db.Repository
	ArchivePath string

	arkv1connect.UnimplementedArkApiHandler
}

func (s *Handler) UploadFile(ctx context.Context, req *connect_go.ClientStream[arkv1.UploadFileRequest]) (*connect_go.Response[arkv1.UploadFileResponse], error) {
	start := time.Now()
	defer func() {
		metrics.UploadFileDurationMs.Observe(float64(time.Since(start).Milliseconds()))
	}()

	next := req.Receive()
	if !next && req.Err() != nil {
		return nil, connect_go.NewError(connect_go.CodeInternal, req.Err())
	}

	metadata := req.Msg().GetMetadata()

	if metadata == nil {
		return nil, connect_go.NewError(connect_go.CodeInvalidArgument, fmt.Errorf("expected metadata"))
	}

	media, err := s.Repo.Get(ctx, metadata.GetHash())
	if err != nil {
		return nil, connect_go.NewError(connect_go.CodeInternal, err)
	}

	if media != nil {
		metrics.TotalDuplicates.Inc()
		return nil, connect_go.NewError(connect_go.CodeAlreadyExists, fmt.Errorf("file already exists: %v", media.Path))
	}

	now := time.Now()
	media = &db.Media{
		Hash:       metadata.GetHash(),
		Path:       metadata.GetName(),
		CreatedAt:  metadata.GetCreatedAt().AsTime(),
		ImportedAt: &now,
	}

	buffer := bytes.Buffer{}
	var size int64

	next = req.Receive()
	for next {
		chunk := req.Msg().GetChunk()

		n, err := buffer.Write(chunk.GetData())
		if err != nil {
			return nil, connect_go.NewError(connect_go.CodeInternal, err)
		}
		size += int64(n)

		next = req.Receive()
	}

	if req.Err() != nil {
		return nil, connect_go.NewError(connect_go.CodeInternal, req.Err())
	}

	if size != metadata.GetSize() {
		return nil, connect_go.NewError(connect_go.CodeInternal, fmt.Errorf("total size mismatch: expected %v, got %v", metadata.GetSize(), size))
	}

	newPath, err := s.copyFile(*media, buffer)
	if err != nil {
		return nil, connect_go.NewError(connect_go.CodeInternal, err)
	}

	media.Path = newPath
	if err := s.Repo.Store(ctx, *media); err != nil {
		return nil, connect_go.NewError(connect_go.CodeInternal, err)
	}

	metrics.TotalImported.Inc()

	return connect_go.NewResponse(&arkv1.UploadFileResponse{}), nil
}

func (s *Handler) copyFile(m db.Media, buffer bytes.Buffer) (string, error) {
	start := time.Now()
	defer func() {
		metrics.CopyFileDurationMs.Observe(float64(time.Since(start).Milliseconds()))
	}()
	year := m.CreatedAt.Format("2006")
	month := m.CreatedAt.Format("01")
	day := m.CreatedAt.Format("02")
	ymdDir := filepath.Join(s.ArchivePath, year, month, day)

	if err := os.MkdirAll(ymdDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("unable to create archive subdirectory %v: %w", ymdDir, err)
	}

	newPath := filepath.Join(ymdDir, filepath.Base(m.Path))
	return newPath, s.atomicallyWriteFile(newPath, bufio.NewReader(&buffer))
}

func (s *Handler) atomicallyWriteFile(filename string, r io.Reader) (err error) {
	dir, file := filepath.Split(filename)
	if dir == "" {
		dir = "."
	}

	tmpDir := filepath.Join(s.ArchivePath, "tmp")
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create temporary subdirectory %v: %w", tmpDir, err)
	}

	f, err := os.CreateTemp(tmpDir, file)
	if err != nil {
		return fmt.Errorf("cannot create temp file: %v", err)
	}
	defer func() {
		if err != nil {
			_ = os.Remove(f.Name())
		}
	}()
	defer f.Close()
	name := f.Name()
	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("cannot write data to tempfile %q: %v", name, err)
	}
	// fsync is important, otherwise os.Rename could rename a zero-length file
	if err := f.Sync(); err != nil {
		return fmt.Errorf("can't flush tempfile %q: %v", name, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("can't close tempfile %q: %v", name, err)
	}

	if err := os.Rename(name, filename); err != nil {
		return fmt.Errorf("cannot replace %q with tempfile %q: %v", filename, name, err)
	}

	return nil
}
