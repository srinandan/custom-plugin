package main

// inspired by https://github.com/salrashid123/envoy_external_authz/blob/master/authz_server/grpc_server.go

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/gogo/googleapis/google/rpc"
)

func getGRPCPort() string {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		return ":50051"
	}
	return port
}

type healthServer struct{}

func (s *healthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	log.Printf("Handling grpc Check request")
	// yeah, right, open 24x7, like 7-11
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

type AuthorizationServer struct{}

func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	log.Println(">>> Authorization called check()")

	b, err := json.MarshalIndent(req.Attributes.Request.Http.Headers, "", "  ")
	if err == nil {
		log.Println("Inbound Headers: ")
		log.Println((string(b)))
	}

	ct, err := json.MarshalIndent(req.Attributes.ContextExtensions, "", "  ")
	if err == nil {
		log.Println("Context Extensions: ")
		log.Println((string(ct)))
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
					{
						Header: &core.HeaderValue{
							Key:   ":host",
							Value: "mocktarget.apigee.net",
						},
						// Note that
						// by Leaving `append` as false, the filter will either add a new header, or override an existing
						// one if there is a match.
						Append: &wrapperspb.BoolValue{Value: false},
					},
				},
			},
		},
	}, nil

}

func main() {

	lis, err := net.Listen("tcp", getGRPCPort())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}
	opts = append(opts)

	s := grpc.NewServer(opts...)

	auth.RegisterAuthorizationServer(s, &AuthorizationServer{})
	healthpb.RegisterHealthServer(s, &healthServer{})

	log.Printf("Starting gRPC Server at %s", getGRPCPort())
	_ = s.Serve(lis)

}
