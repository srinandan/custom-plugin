package extauthz

import (
  "encoding/json"
  "log"

  "golang.org/x/net/context"
  "google.golang.org/grpc"
  //"google.golang.org/protobuf/types/known/wrapperspb"

  "google.golang.org/grpc/codes"
  healthpb "google.golang.org/grpc/health/grpc_health_v1"
  "google.golang.org/grpc/status"

  rpcstatus "google.golang.org/genproto/googleapis/rpc/status"

  core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
  auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
  "github.com/gogo/googleapis/google/rpc"
)

type healthServer struct{}

func (s *healthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

// Register registers
func (a *AuthorizationServer) Register(s *grpc.Server) {
	auth.RegisterAuthorizationServer(s, a)
}

// AuthorizationServer server
type AuthorizationServer struct {}

func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	log.Println(">>> Authorization called check()")

	if req.Attributes != nil &&
		req.Attributes.Request != nil &&
		req.Attributes.Request.Http != nil &&
		req.Attributes.Request.Http.Headers != nil {
		b, err := json.MarshalIndent(req.Attributes.Request.Http.Headers, "", "  ")
		if err == nil {
			log.Println("Inbound Headers: ")
			log.Println((string(b)))
		}
	}

	if req.Attributes != nil && req.Attributes.ContextExtensions != nil {
		ct, err := json.MarshalIndent(req.Attributes.ContextExtensions, "", "  ")
		if err == nil {
			log.Println("Context Extensions: ")
			log.Println((string(ct)))
		}
	}

	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(rpc.OK),
		},
		HttpResponse: &auth.CheckResponse_OkResponse{
			OkResponse: &auth.OkHttpResponse{
				Headers: []*core.HeaderValueOption{
					{
						Header: &core.HeaderValue{
							Key:   "x-custom-header-from-authz",
							Value: "some value",
						},
					},
					/*{
						Header: &core.HeaderValue{
							Key:   "host",
							Value: "mocktarget.apigee.net",
						},
						// Note that
						// by Leaving `append` as false, the filter will either add a new header, or override an existing
						// one if there is a match.
						//Append: &wrapperspb.BoolValue{Value: false},
					},*/
				},
			},
		},
	}, nil

}
