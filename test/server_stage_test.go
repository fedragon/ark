package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	arkv1 "github.com/fedragon/ark/gen/ark/v1"
	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/db"
	"github.com/fedragon/ark/internal/fs"
	"github.com/fedragon/ark/internal/server"
	_ "github.com/fedragon/ark/testing"

	"github.com/bufbuild/connect-go"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServerStage struct {
	t           *testing.T
	server      *httptest.Server
	client      arkv1connect.ArkApiClient
	uploadError error
}

func NewServerStage(t *testing.T) *ServerStage {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
	})
	repo := db.NewRedisRepository(client)

	handler := &server.Handler{
		Repo:        repo,
		ArchivePath: "./archive",
	}

	mux := http.NewServeMux()
	mux.Handle(arkv1connect.NewArkApiHandler(handler))

	us := httptest.NewUnstartedServer(mux)
	us.EnableHTTP2 = true
	us.Start()

	return &ServerStage{
		t:      t,
		server: us,
		client: arkv1connect.NewArkApiClient(http.DefaultClient, us.URL),
	}
}

func (s *ServerStage) And() *ServerStage {
	return s
}

func (s *ServerStage) Given() *ServerStage {
	return s
}

func (s *ServerStage) When() *ServerStage {
	return s
}

func (s *ServerStage) Then() *ServerStage {
	return s
}

func (s *ServerStage) FileDoesNotExist() *ServerStage {
	return s
}

func (s *ServerStage) FileExists() *ServerStage {
	return s
}

func (s *ServerStage) ClientUploadsFile(path string) *ServerStage {
	data, err := os.ReadFile(path)
	assert.NoError(s.t, err)

	hash, err := fs.Hash(path)
	assert.NoError(s.t, err)

	stream := s.client.UploadFile(context.Background())

	err = stream.Send(&arkv1.UploadFileRequest{
		File: &arkv1.UploadFileRequest_Metadata{
			Metadata: &arkv1.Metadata{
				Hash:      hash,
				Name:      path,
				Size:      int64(len(data)),
				CreatedAt: timestamppb.New(time.Now()),
			},
		},
	})
	assert.NoError(s.t, err)

	err = stream.Send(&arkv1.UploadFileRequest{
		File: &arkv1.UploadFileRequest_Chunk{
			Chunk: &arkv1.Chunk{
				Data: data,
			},
		},
	})
	assert.NoError(s.t, err)
	_, s.uploadError = stream.CloseAndReceive()

	return s
}

func (s *ServerStage) ClientUploadsFileAgain(path string) *ServerStage {
	data, err := os.ReadFile(path)
	assert.NoError(s.t, err)

	hash, err := fs.Hash(path)
	assert.NoError(s.t, err)

	stream := s.client.UploadFile(context.Background())
	err = stream.Send(&arkv1.UploadFileRequest{
		File: &arkv1.UploadFileRequest_Metadata{
			Metadata: &arkv1.Metadata{
				Hash:      hash,
				Name:      path,
				Size:      int64(len(data)),
				CreatedAt: timestamppb.New(time.Now()),
			},
		},
	})
	assert.NoError(s.t, err)
	_, s.uploadError = stream.CloseAndReceive()

	return s
}

func (s *ServerStage) UploadSucceeds() *ServerStage {
	assert.NoError(s.t, s.uploadError)
	return s
}

func (s *ServerStage) UploadIsSkipped() *ServerStage {
	target := &connect.Error{}
	if assert.Error(s.t, s.uploadError) && assert.ErrorAs(s.t, s.uploadError, &target) {
		assert.Equal(s.t, connect.CodeAlreadyExists, target.Code(), target.Error())
	} else {
		s.t.FailNow()
	}

	return s
}
