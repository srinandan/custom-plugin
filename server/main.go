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

import (
	"log"
	"net"
	"os"
	"os/signal"

	extauthz "github.com/srinandan/custom-plugin/server/extauthz"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func getGRPCPort() string {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		return ":5000"
	}
	return port
}

func main() {

	ctx := context.Background()

	listen, err := net.Listen("tcp", getGRPCPort())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}

	grpcServer := grpc.NewServer(opts...)

	as := &extauthz.AuthorizationServer{}
	as.Register(grpcServer)

	//healthpb.RegisterHealthServer(server, &healthServer{})

	log.Printf("starting gRPC Server at %s", getGRPCPort())

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Printf("shutting down gRPC server...")

			grpcServer.GracefulStop()

			<-ctx.Done()
		}
	}()

	_ = grpcServer.Serve(listen)

}
