// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// inspired by https://github.com/salrashid123/envoy_external_authz/blob/master/authz_server/grpc_server.go

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/gogo/googleapis/google/rpc"
)

func getGRPCPort() string {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		return "0.0.0.0:50051"
	}
	return port
}

type healthServer struct{}

func (s *healthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

type AuthorizationServer struct{}

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
					{
						Header: &core.HeaderValue{
							Key:   "host",
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

	ctx := context.Background()

	listen, err := net.Listen("tcp", getGRPCPort())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}

	server := grpc.NewServer(opts...)

	auth.RegisterAuthorizationServer(server, &AuthorizationServer{})
	healthpb.RegisterHealthServer(server, &healthServer{})

	log.Printf("starting gRPC Server at %s", getGRPCPort())

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Printf("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	_ = server.Serve(listen)

}
