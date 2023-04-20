package importer

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	arkv1 "github.com/fedragon/ark/gen/ark/v1"
	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/fs"

	connect "github.com/bufbuild/connect-go"
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
}

func (imp *Imp) Import(ctx context.Context, sourceDir string) error {
	group := errgroup.Group{}
	sendOne := func(ctx context.Context, in <-chan db.Media) error {
		for m := range in {
			_, err := imp.sendMedia(ctx, m)

			var cerr *connect.Error
			if err != nil {
				if !errors.As(err, &cerr) {
					return err
				}

				if cerr.Code() == connect.CodeAlreadyExists {
					fmt.Printf("skipped duplicate %s\n", m.Path)
					continue
				}
			}

			fmt.Printf("imported %s\n", m.Path)
		}

		return nil
	}

	allMedia := fs.Walk(sourceDir, imp.FileTypes)
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		group.Go(func() error { return sendOne(ctx, allMedia) })
	}
	//
	//for m := range media {
	//	if _, err := imp.sendMedia(ctx, m); err != nil {
	//		var cerr *connect.Error
	//		if errors.As(err, &cerr) {
	//			if cerr.Code() == connect.CodeAlreadyExists {
	//				fmt.Printf("skipped duplicate %s\n", m.Path)
	//				continue
	//			}
	//		}
	//		return err
	//	}
	//
	//	fmt.Printf("imported %s\n", m.Path)
	//}
	return group.Wait()
}

func (imp *Imp) sendMedia(ctx context.Context, m db.Media) (*connect.Response[arkv1.UploadFileResponse], error) {
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
