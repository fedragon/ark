package test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	arkv1 "github.com/fedragon/ark/gen/ark/v1"
	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/importer"
	_ "github.com/fedragon/ark/testing"

	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockArkApiServer struct {
	FileTypes []string

	arkv1connect.UnimplementedArkApiHandler

	uploadFileResponse *arkv1.UploadFileResponse
	uploadFileError    error
}

func (maas *MockArkApiServer) UploadFile(_ context.Context, _ *connect.ClientStream[arkv1.UploadFileRequest]) (*connect.Response[arkv1.UploadFileResponse], error) {
	return connect.NewResponse(maas.uploadFileResponse), maas.uploadFileError
}

func (maas *MockArkApiServer) setUploadFileResponse(response *arkv1.UploadFileResponse) {
	maas.uploadFileResponse = response
}

func (maas *MockArkApiServer) setUploadFileError(err error) {
	maas.uploadFileError = err
}

type ClientStage struct {
	t           *testing.T
	imp         importer.Importer
	mock        *MockArkApiServer
	server      *httptest.Server
	importError error
}

func NewClientStage(t *testing.T) *ClientStage {
	types := []string{"jpg"}
	mock := &MockArkApiServer{FileTypes: types}

	mux := http.NewServeMux()
	mux.Handle(arkv1connect.NewArkApiHandler(mock))

	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.Start()

	return &ClientStage{
		t: t,
		imp: &importer.Imp{
			Client: arkv1connect.NewArkApiClient(
				http.DefaultClient,
				server.URL,
				connect.WithSendGzip(),
			),
			FileTypes: types,
			Logger:    zap.NewNop(),
		},
		mock:   mock,
		server: server,
	}
}

func (s *ClientStage) And() *ClientStage {
	return s
}

func (s *ClientStage) Given() *ClientStage {
	return s
}

func (s *ClientStage) When() *ClientStage {
	return s
}

func (s *ClientStage) Then() *ClientStage {
	return s
}

func (s *ClientStage) UploadFileWillSucceed() *ClientStage {
	s.mock.setUploadFileResponse(&arkv1.UploadFileResponse{})
	return s
}

func (s *ClientStage) UploadFileWillBeSkipped() *ClientStage {
	s.mock.setUploadFileError(connect.NewError(connect.CodeAlreadyExists, errors.New("file already exists")))
	return s
}

func (s *ClientStage) UploadFileWillFail() *ClientStage {
	s.mock.setUploadFileError(connect.NewError(connect.CodeInternal, errors.New("something went wrong")))
	return s
}

func (s *ClientStage) ClientUploadsFile() *ClientStage {
	s.importError = s.imp.Import(context.Background(), "./test/data/doge.jpg")
	return s
}

func (s *ClientStage) ImportSucceeds() *ClientStage {
	assert.NoError(s.t, s.importError)
	return s
}
