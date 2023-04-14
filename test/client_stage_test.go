package test

import (
	"context"
	"errors"
	arkv1 "github.com/fedragon/ark/gen/ark/v1"
	"github.com/fedragon/ark/gen/ark/v1/arkv1connect"
	"github.com/fedragon/ark/internal/importer"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	connect_go "github.com/bufbuild/connect-go"
)

type MockArkApiServer struct {
	FileTypes []string

	arkv1connect.UnimplementedArkApiHandler

	fileExistsResponse *arkv1.FileExistsResponse
	fileExistsError    error
	uploadFileResponse *arkv1.UploadFileResponse
	uploadFileError    error
}

func (maas *MockArkApiServer) FileExists(context.Context, *connect_go.Request[arkv1.FileExistsRequest]) (*connect_go.Response[arkv1.FileExistsResponse], error) {
	return connect_go.NewResponse(maas.fileExistsResponse), maas.fileExistsError
}

func (maas *MockArkApiServer) UploadFile(_ context.Context, stream *connect_go.ClientStream[arkv1.UploadFileRequest]) (*connect_go.Response[arkv1.UploadFileResponse], error) {
	for stream.Receive() {
		// consume the stream
	}

	return connect_go.NewResponse(maas.uploadFileResponse), maas.uploadFileError
}

func (maas *MockArkApiServer) setFileExistsResponse(response *arkv1.FileExistsResponse) {
	maas.fileExistsResponse = response
}

func (maas *MockArkApiServer) setFileExistsError(err error) {
	maas.fileExistsError = err
}

func (maas *MockArkApiServer) setUploadFileResponse(response *arkv1.UploadFileResponse) {
	maas.uploadFileResponse = response
}

func (maas *MockArkApiServer) setUploadFileError(err error) {
	maas.uploadFileError = err
}

type ClientStage struct {
	t            *testing.T
	imp          importer.Importer
	mock         *MockArkApiServer
	server       *httptest.Server
	importResult error
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
			Client:    arkv1connect.NewArkApiClient(http.DefaultClient, server.URL),
			FileTypes: types,
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

func (s *ClientStage) FileDoesNotExist() *ClientStage {
	s.mock.setFileExistsResponse(&arkv1.FileExistsResponse{Exists: false})
	return s
}

func (s *ClientStage) FileExists() *ClientStage {
	s.mock.setFileExistsResponse(&arkv1.FileExistsResponse{Exists: true})
	return s
}

func (s *ClientStage) UploadFileWillSucceed() *ClientStage {
	s.mock.setUploadFileResponse(&arkv1.UploadFileResponse{Success: true})
	return s
}

func (s *ClientStage) UploadFileWillFail() *ClientStage {
	s.mock.setUploadFileError(errors.New("upload failed"))
	return s
}

func (s *ClientStage) ClientUploadsFile() *ClientStage {
	s.importResult = s.imp.Import(context.Background(), "./data/doge.jpg")
	return s
}

func (s *ClientStage) ImportSucceeds() *ClientStage {
	assert.NoError(s.t, s.importResult)
	return s
}

func (s *ClientStage) ImportFails() *ClientStage {
	assert.Error(s.t, s.importResult)
	return s
}
