package main

import (
	"google.golang.org/grpc"
	"net"
	"utils"
)

type server struct {
	grpcServer *grpc.Server
	Logger     *utils.Logger
}

func (s *server) start(address string)  {
	s.Logger.InfoPrintln("Starting server...")
	// start listener
	listener, err := net.Listen("tcp", address)
	if err != nil {
		s.Logger.ErrorLogger.Fatalf("Could not listen at ip address %v. :: %v", address, err)
	}

	// Accept incoming connections on the listener
	s.Logger.InfoPrintln("Server started.")
	err = s.grpcServer.Serve(listener)
	if err != nil {
		s.Logger.ErrorLogger.Fatalf("Failed to start gRPC server. :: %v", err)
	}
}

func newServer(logger *utils.Logger) *server {
	return &server{
		grpcServer: grpc.NewServer(),
		Logger:     logger,
	}
}

