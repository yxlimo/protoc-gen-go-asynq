package gen

const fednodeTpl = `package {{ package . }}
import (
	"context"
	"fmt"

	"git.basebit.me/xee/fednode/proto/basebit/fednodex/v2"
	"github.com/cockroachdb/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	"git.basebit.me/enigma/toolx/efednode"
	"git.basebit.me/enigma/toolx/efednode/proto"
)

{{ range .Services }}

type {{ .Name }}FednodeClient interface {
	WithNode(nodeID string) {{ .Name }}Client
}

func New{{ .Name }}FednodeClient(client *efednode.Module) {{ .Name }}FednodeClient {
	return &fednode{{ .Name }}Client{
		client: client,
	}
}

type fednode{{ .Name }}Client struct {
	client *efednode.Module
}

func (f *fednode{{ .Name }}Client) WithNode(nodeID string) {{ .Name }}Client {
	return &fednode{{ .Name }}NodeClient{
		client: f.client.Node(nodeID),
	}
}

type fednode{{ .Name }}NodeClient struct {
	client *efednode.Node
}

{{ range .Methods }}
{{ if and .ServerStreaming .ClientStreaming }}
func (f *fednode{{ .Service.Name }}NodeClient) {{ name . }}(ctx context.Context, opts ...grpc.CallOption) ({{ .Service.Name}}_{{ name . }}Client, error) {
	return nil, status.Error(codes.Unimplemented, "client stream method is not supported")
}
{{ else if .ClientStreaming}}
func (f *fednode{{ .Service.Name }}NodeClient) {{ name . }}(ctx context.Context, opts ...grpc.CallOption) ({{ .Service.Name}}_{{ name . }}Client, error) {
	return nil, status.Error(codes.Unimplemented, "client stream method is not supported")
}
{{ else if .ServerStreaming}}
func (f *fednode{{ .Service.Name }}NodeClient) {{ name . }}(ctx context.Context, in *{{ name .Input}}, opts ...grpc.CallOption) ({{ .Service.Name}}_{{ name . }}Client, error) {
	return nil, status.Error(codes.Unimplemented, "client stream method is not supported")
}
{{ else }}
func (f *fednode{{ .Service.Name }}NodeClient) {{ name . }}(ctx context.Context, in *{{ name .Input}}, opts ...grpc.CallOption) (*{{ name .Output}}, error) {
	payload, err := anypb.New(in)
	if err != nil {
		return nil, errors.Wrap(err, "create payload")
	}

	msg := &proto.MessageRequest{
		Method:  "{{ .Service.Name }}.{{ name . }}",
		Payload: payload,
	}
	res, err := f.client.Call(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "call")
	}
	if status.FromProto(res.GetStatus()).Code() != codes.OK {
		return nil, status.FromProto(res.GetStatus()).Err()
	}
	pbRes := &{{ name .Output}}{}
	if err := res.Result.UnmarshalTo(pbRes); err != nil {
		return nil, fmt.Errorf("unmarshal proto(%+v): %w", pbRes, err)
	}
	return pbRes, nil
}
{{ end }}
{{ end }}


type handler{{ .Name }} struct {
	fednodex.UnimplementedFedServiceServer
	srv {{ .Name }}Server
}

func (h *handler{{ .Name }}) OnFedMessage(ctx context.Context, fedReq *fednodex.OnFedMessageRequest) (*fednodex.OnFedMessageResponse, error) {
	msgReq := &proto.MessageRequest{}
	if err := fedReq.GetMessage().UnmarshalTo(msgReq); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	switch msgReq.GetMethod() {
	{{ range .Methods }}
	{{ if or .ClientStreaming .ServerStreaming }}
	// skip steaming {{ .Service.Name }}.{{ .Name }}
	{{ else }}
	case "{{ .Service.Name }}.{{ name . }}":
		req := &{{ name .Input }}{}
		if err := msgReq.Payload.UnmarshalTo(req); err != nil {
			return efednode.ResponseError(err), nil
		}
		res, err := h.srv.{{ name . }}(ctx, req)
		if err != nil {
			return efednode.ResponseError(err), nil
		}
		return efednode.ResponseOK(res), nil
	{{ end }}
	{{ end }}
	}
	return efednode.ResponseError(status.Error(codes.Unimplemented, "unknown method")), nil
}

func Register{{ .Name }}FedNodeServer(s *grpc.Server, srv {{ .Name }}Server) {
	fednodex.RegisterFedServiceServer(s, &handler{{ .Name }}{srv: srv})
}

{{ end }}
`
