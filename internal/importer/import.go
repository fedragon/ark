package importer

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"runtime"

	arkv1 "github.com/fedragon/ark/gen/ark/v1"
	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/fs"

	"connectrpc.com/connect"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Importer interface {
	// Import imports all files in sourceDir, skipping duplicates
	Import(ctx context.Context, sourceDir string) error
}

type Imp struct {
	Client    arkv1connect.ArkApiClient
	FileTypes []string
	Logger    *zap.Logger
}

func (imp *Imp) Import(ctx context.Context, sourceDir string) error {
	group := errgroup.Group{}
	sendOne := func(ctx context.Context, in <-chan db.Media) error {
		for m := range in {
			if m.Err != nil {
				return m.Err
			}

			if _, err := imp.send(ctx, m); err != nil {
				var cerr *connect.Error
				if !errors.As(err, &cerr) || cerr.Code() != connect.CodeAlreadyExists {
					return err
				}

				imp.Logger.Info("Skipped duplicate file %s", zap.String("path", m.Path))
				continue
			}

			imp.Logger.Info("Imported file", zap.String("path", m.Path))
		}

		return nil
	}

	allMedia := fs.Walk(sourceDir, imp.FileTypes)
	for i := 0; i < runtime.NumCPU(); i++ {
		group.Go(func() error { return sendOne(ctx, allMedia) })
	}

	return group.Wait()
}

func (imp *Imp) send(ctx context.Context, m db.Media) (*connect.Response[arkv1.UploadFileResponse], error) {
	file, err := os.Open(m.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	stream := imp.Client.UploadFile(ctx)
	err = stream.Send(&arkv1.UploadFileRequest{
		File: &arkv1.UploadFileRequest_Metadata{
			Metadata: &arkv1.Metadata{
				Hash:      m.Hash,
				Name:      m.Path,
				Size:      stat.Size(),
				CreatedAt: timestamppb.New(stat.ModTime()),
			},
		},
	})
	if err != nil {
		if errors.Is(err, io.EOF) {
			return stream.CloseAndReceive()
		}
		return nil, err
	}

	reader := bufio.NewReader(file)
	chunk := make([]byte, 1024*1024)

	for {
		n, err := reader.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		err = stream.Send(&arkv1.UploadFileRequest{
			File: &arkv1.UploadFileRequest_Chunk{
				Chunk: &arkv1.Chunk{
					Data: chunk[:n],
				},
			},
		})

		if err != nil {
			if errors.Is(err, io.EOF) {
				return stream.CloseAndReceive()
			}
		}
	}

	return stream.CloseAndReceive()
}
