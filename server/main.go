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
	"syscall"
	"time"

	extauthz "github.com/srinandan/custom-plugin/server/extauthz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"github.com/srinandan/custom-plugin/routes"
)

func getGRPCPort() string {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		return ":5000"
	}
	return port
}

func main() {
	if err := routes.ReadRoutesFile(); err != nil {
		fmt.Errorf("unable to load routing table.")
		os.Exit(1)
	}
	serve()
	select {}
}

func serve() {
	// gRPC server
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: 10 * time.Minute,
		}),
		grpc.MaxConcurrentStreams(10),
	}

	grpcServer := grpc.NewServer(opts...)

	as := &extauthz.AuthorizationServer{}
	as.Register(grpcServer)

	// grpc health
	grpcHealth := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, grpcHealth)

	log.Printf("starting gRPC Server at %s", getGRPCPort())

	// grpc listener
	grpcListener, err := net.Listen("tcp", getGRPCPort())
	if err != nil {
		panic(err)
	}

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("%s", err)
		}
	}()

	// watch for termination signals
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)    // terminal
		signal.Notify(sigint, syscall.SIGTERM) // kubernetes
		sig := <-sigint
		log.Printf("shutdown signal: %s", sig)
		signal.Stop(sigint)

		grpcServer.GracefulStop()

		log.Println("shutdown complete")
		os.Exit(0)
	}()
}
