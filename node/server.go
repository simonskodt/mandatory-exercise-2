package main

import (
	"google.golang.org/grpc"
	"net"
	"service"
	"strconv"
	"utils"
)

type server struct {
	grpcServer *grpc.Server
	port   int
	logger *utils.Logger
}

func (s *server) start(address string, port int) {
	s.logger.InfoPrintln("Starting server...")
	// start listener
	listener, err := net.Listen("tcp", address+":"+strconv.Itoa(port))
	if err != nil {
		s.logger.ErrorFatalf("Could not listen at ip address %v. :: %v", address, err)
	}

	// Register the gRPC server on the node service
	service.RegisterServiceServer(s.grpcServer, n)

	// Accept incoming connections on the listener
	s.logger.InfoPrintln("Server started.")
	err = s.grpcServer.Serve(listener)
	if err != nil {
		s.logger.ErrorFatalf("Failed to start gRPC server. :: %v", err)
	}
}

func (s *server) stop() {
	s.logger.WarningPrintln("Stopping server...")

}

func newServer(port int, logger *utils.Logger) *server {
	return &server{
		grpcServer: grpc.NewServer(),
		port:   port,
		logger: logger,
	}
}
