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
	"github.com/fedragon/ark/internal/server"
	_ "github.com/fedragon/ark/testing"

	connect_go "github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var hash = []byte("38bcf7453c666c5a901c7a7547fca31167d8d736aa4581a44e431772deeae73fdfc0b76962d42bd6565bb249c92cab86166bd19f06f04096573f835b52a8e2fdd59a00791e55b607473ac8aa5f8f624b3ca63a733d6ee704f309e374e055dd25001bd65a0421ebf1029bd7673787964a80d99743caae94fad54d8f550ce79e8a4a7e93b00f140be195140ef9f6814294f6d6b4df4519c9d7845902f26146cdaee1fc5e763ef7e3d2b88cd39c726a807323a03b53f37a97df2e60c29ede044ca0252ff7c533a3e467044028639d547371ff2e3019e66272aed73cb61b982fefab7cb51e3a908e7add8086b25366366261885027bae0a166965179f88a7eaf8ac9")

type ServerStage struct {
	t           *testing.T
	server      *httptest.Server
	client      arkv1connect.ArkApiClient
	uploadError error
}

func NewServerStage(t *testing.T) *ServerStage {
	repo, err := db.NewSqlite3Repository("./ark.db")
	if err != nil {
		t.Fatal(err.Error())
	}

	handler := &server.Ark{
		Repo:        repo,
		FileTypes:   []string{"jpg"},
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

func (s *ServerStage) ClientUploadsFile() *ServerStage {
	data, err := os.ReadFile("./test/data/doge.jpg")
	assert.NoError(s.t, err)

	stream := s.client.UploadFile(context.Background())

	err = stream.Send(&arkv1.UploadFileRequest{
		File: &arkv1.UploadFileRequest_Metadata{
			Metadata: &arkv1.Metadata{
				Hash:      hash,
				Name:      "./data/doge.jpg",
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
				Size: int64(len(data)),
			},
		},
	})
	assert.NoError(s.t, err)

	_, err = stream.CloseAndReceive()
	s.uploadError = err

	return s
}

func (s *ServerStage) UploadSucceeds() *ServerStage {
	assert.NoError(s.t, s.uploadError)
	return s
}

func (s *ServerStage) UploadIsSkipped() *ServerStage {
	target := &connect_go.Error{}
	if assert.Error(s.t, s.uploadError) && assert.ErrorAs(s.t, s.uploadError, &target) {
		assert.Equal(s.t, connect_go.CodeAlreadyExists, target.Code(), target.Error())
	}

	return s
}

func (s *ServerStage) UploadFails() *ServerStage {
	target := &connect_go.Error{}
	if assert.Error(s.t, s.uploadError) && assert.ErrorAs(s.t, s.uploadError, &target) {
		assert.NotEqual(s.t, connect_go.CodeAlreadyExists, target.Code(), target.Error())
	}

	return s
}
