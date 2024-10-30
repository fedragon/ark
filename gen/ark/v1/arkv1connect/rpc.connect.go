// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: ark/v1/rpc.proto

package arkv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/fedragon/ark/gen/ark/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// ArkApiName is the fully-qualified name of the ArkApi service.
	ArkApiName = "ark.v1.ArkApi"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ArkApiUploadFileProcedure is the fully-qualified name of the ArkApi's UploadFile RPC.
	ArkApiUploadFileProcedure = "/ark.v1.ArkApi/UploadFile"
)

// ArkApiClient is a client for the ark.v1.ArkApi service.
type ArkApiClient interface {
	UploadFile(context.Context) *connect_go.ClientStreamForClient[v1.UploadFileRequest, v1.UploadFileResponse]
}

// NewArkApiClient constructs a client for the ark.v1.ArkApi service. By default, it uses the
// Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewArkApiClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ArkApiClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &arkApiClient{
		uploadFile: connect_go.NewClient[v1.UploadFileRequest, v1.UploadFileResponse](
			httpClient,
			baseURL+ArkApiUploadFileProcedure,
			opts...,
		),
	}
}

// arkApiClient implements ArkApiClient.
type arkApiClient struct {
	uploadFile *connect_go.Client[v1.UploadFileRequest, v1.UploadFileResponse]
}

// UploadFile calls ark.v1.ArkApi.UploadFile.
func (c *arkApiClient) UploadFile(ctx context.Context) *connect_go.ClientStreamForClient[v1.UploadFileRequest, v1.UploadFileResponse] {
	return c.uploadFile.CallClientStream(ctx)
}

// ArkApiHandler is an implementation of the ark.v1.ArkApi service.
type ArkApiHandler interface {
	UploadFile(context.Context, *connect_go.ClientStream[v1.UploadFileRequest]) (*connect_go.Response[v1.UploadFileResponse], error)
}

// NewArkApiHandler builds an HTTP handler from the service implementation. It returns the path on
// which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewArkApiHandler(svc ArkApiHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	arkApiUploadFileHandler := connect_go.NewClientStreamHandler(
		ArkApiUploadFileProcedure,
		svc.UploadFile,
		opts...,
	)
	return "/ark.v1.ArkApi/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ArkApiUploadFileProcedure:
			arkApiUploadFileHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedArkApiHandler returns CodeUnimplemented from all methods.
type UnimplementedArkApiHandler struct{}

func (UnimplementedArkApiHandler) UploadFile(context.Context, *connect_go.ClientStream[v1.UploadFileRequest]) (*connect_go.Response[v1.UploadFileResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("ark.v1.ArkApi.UploadFile is not implemented"))
}
