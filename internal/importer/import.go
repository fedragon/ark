package importer

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	arkv1 "github.com/fedragon/ark/gen/ark/v1"
	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/fs"

	connect_go "github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Importer interface {
	// Import imports all files in sourceDir, skipping duplicates
	Import(ctx context.Context, sourceDir string) error
}

type Imp struct {
	Client    arkv1connect.ArkApiClient
	FileTypes []string
}

func (imp *Imp) Import(ctx context.Context, sourceDir string) error {
	for m := range fs.Walk(sourceDir, imp.FileTypes) {
		if _, err := imp.sendMedia(ctx, m); err != nil {
			var cerr *connect_go.Error
			if errors.As(err, &cerr) {
				if cerr.Code() == connect_go.CodeAlreadyExists {
					fmt.Printf("skipped duplicate %s\n", m.Path)
					continue
				}
			}
			return err
		}

		fmt.Printf("imported %s\n", m.Path)
	}

	return nil
}

func (imp *Imp) sendMedia(ctx context.Context, m db.Media) (*connect_go.Response[arkv1.UploadFileResponse], error) {
	file, err := os.Open(m.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := os.Stat(m.Path)
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
				CreatedAt: timestamppb.New(time.Now()),
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
