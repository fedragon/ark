package server

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	arkv1 "github.com/fedragon/ark/gen/ark/v1"
	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/image"
	"github.com/fedragon/ark/internal/metrics"

	"connectrpc.com/connect"
)

type Handler struct {
	Repo        db.Repository
	ArchivePath string

	arkv1connect.UnimplementedArkApiHandler
}

func (s *Handler) UploadFile(ctx context.Context, req *connect.ClientStream[arkv1.UploadFileRequest]) (*connect.Response[arkv1.UploadFileResponse], error) {
	start := time.Now()
	defer func() {
		metrics.UploadFileDurationMs.Observe(float64(time.Since(start).Milliseconds()))
	}()

	next := req.Receive()
	if !next && req.Err() != nil {
		return nil, connect.NewError(connect.CodeInternal, req.Err())
	}

	metadata := req.Msg().GetMetadata()
	if metadata == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("expected metadata"))
	}

	media, err := s.Repo.Get(ctx, metadata.GetHash())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if media != nil {
		metrics.TotalDuplicates.Inc()
		return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("file already exists: %v", media.Path))
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
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		size += int64(n)

		next = req.Receive()
	}

	if req.Err() != nil {
		return nil, connect.NewError(connect.CodeInternal, req.Err())
	}

	if size != metadata.GetSize() {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("total size mismatch: expected %v, got %v", metadata.GetSize(), size))
	}

	newPath, err := s.copyFile(*media, buffer)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	media.Path = newPath
	if err := s.Repo.Store(ctx, *media); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	metrics.TotalImported.Inc()

	return connect.NewResponse(&arkv1.UploadFileResponse{}), nil
}

func (s *Handler) copyFile(m db.Media, buffer bytes.Buffer) (string, error) {
	start := time.Now()
	defer func() {
		metrics.CopyFileDurationMs.Observe(float64(time.Since(start).Milliseconds()))
	}()

	filename := filepath.Base(m.Path)
	tmpPath, err := s.writeFile(filename, bufio.NewReader(&buffer))
	if err != nil {
		return "", fmt.Errorf("unable to write temp file: %w", err)
	}

	createdAt, err := image.ParseCreatedAt(tmpPath)
	if err != nil {
		if !errors.Is(err, &image.ErrNotFound{}) {
			return "", fmt.Errorf("unable to parse createdAt: %w", err)
		}

		createdAt = m.CreatedAt
	}

	year := createdAt.Format("2006")
	month := createdAt.Format("01")
	day := createdAt.Format("02")
	ymdDir := filepath.Join(s.ArchivePath, year, month, day)

	if err := os.MkdirAll(ymdDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("unable to create archive subdirectory %v: %w", ymdDir, err)
	}

	newPath := filepath.Join(ymdDir, filename)
	if err := os.Rename(tmpPath, newPath); err != nil {
		return "", fmt.Errorf("cannot replace %s with temp file %s: %v", tmpPath, newPath, err)
	}

	return newPath, nil
}

func (s *Handler) writeFile(filename string, r io.Reader) (string, error) {
	tmpDir := filepath.Join(s.ArchivePath, "tmp")
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("unable to create temporary subdirectory %v: %w", tmpDir, err)
	}

	ext := filepath.Ext(filename)
	f, err := os.CreateTemp(tmpDir, fmt.Sprintf("%s.*%s", filename, ext))
	if err != nil {
		return "", fmt.Errorf("cannot create temp file: %v", err)
	}
	defer func() {
		if err != nil {
			_ = os.Remove(f.Name())
		}
	}()
	defer f.Close()
	name := f.Name()
	if _, err := io.Copy(f, r); err != nil {
		return "", fmt.Errorf("cannot write data to temp file %q: %v", name, err)
	}
	// fsync is important, otherwise os.Rename could rename a zero-length file
	if err := f.Sync(); err != nil {
		return "", fmt.Errorf("cannot flush temp file %q: %v", name, err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("cannot close temp file %q: %v", name, err)
	}

	return name, nil
}
