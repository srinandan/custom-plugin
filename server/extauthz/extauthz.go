package extauthz

import (
	"encoding/json"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"

	corev2 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/gogo/googleapis/google/rpc"
)

// inspired by https://github.com/salrashid123/envoy_external_authz/blob/master/authz_server/grpc_server.go

// Register registers
func (a *AuthorizationServer) Register(s *grpc.Server) {
	auth.RegisterAuthorizationServer(s, a)
}

// AuthorizationServer server
type AuthorizationServer struct{}

func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	log.Println(">>> Authorization called check()")
	backend := "default"
	basePath := "/"

	if req.Attributes != nil &&
		req.Attributes.Request != nil &&
		req.Attributes.Request.Http != nil &&
		req.Attributes.Request.Http.Headers != nil {
		if b, err := json.MarshalIndent(req.Attributes.Request.Http.Headers, "", "  "); err == nil {
			log.Println("Inbound Headers: ")
			log.Println((string(b)))
			switch backend = req.Attributes.Request.Http.Headers["x-backend-url"]; backend {
			case "mocktarget":
				backend = "mocktarget.apigee.net"
				basePath = "/iloveapis"
			case "postman":
				backend = "postman-echo.com"
				basePath = "/postman"
			default:
				backend = "default"
				basePath = "/"
			}
		}
	}

	if req.Attributes != nil &&
		req.Attributes.Request != nil &&
		req.Attributes.Request.Http != nil &&
		req.Attributes.Request.Http.Body != "" {

		log.Println("Payload >> ", req.Attributes.Request.Http.Body)
	}

	if req.Attributes != nil && req.Attributes.ContextExtensions != nil {
		if ct, err := json.MarshalIndent(req.Attributes.ContextExtensions, "", "  "); err == nil {
			log.Println("Context Extensions: ")
			log.Println((string(ct)))
		}
	}

	return checkResponse(backend, basePath), nil
}

func checkResponse(backend string, basePath string) (*auth.CheckResponse) {
	log.Println(">>> Authorization CheckResponse_OkResponse")
	log.Println("Selecting route ", backend)

	if backend == "default" {
		return &auth.CheckResponse{
			Status: &rpcstatus.Status{
				Code: int32(rpc.OK),
			},
			HttpResponse: &auth.CheckResponse_OkResponse{
				OkResponse: &auth.OkHttpResponse{},
			},
		}
	} else {
		return &auth.CheckResponse{
			Status: &rpcstatus.Status{
				Code: int32(rpc.OK),
			},
			HttpResponse: &auth.CheckResponse_OkResponse{
				OkResponse: &auth.OkHttpResponse{
					Headers: []*corev2.HeaderValueOption{
						setHeader("host", backend, false),
						setHeader(":path", basePath, false),
					},
				},
			},
		}
	}
}

func setHeader(name string, value string, append bool) *corev2.HeaderValueOption {
	header := &corev2.HeaderValue{
		Key:   name,
		Value: value,
	}

	return &corev2.HeaderValueOption{
		Header: header,
		Append: &wrapperspb.BoolValue{Value: append},
	}
}
