package server

import (
	"google.golang.org/grpc"
	"mandatory-exercise-2/service"
	"mandatory-exercise-2/utils"
	"net"
)

type Server struct {
	service.UnimplementedServiceServer
	ipAddress string
	logger *utils.Logger
}

func (s *Server) StartServer()  {
	s.logger.InfoPrintln("Starting server...")
	// Init listener
	listener, err := net.Listen("tcp", s.ipAddress)
	if err != nil {
		s.logger.ErrorLogger.Fatalf("Could not listen at ip address %v. :: %v", s.ipAddress, err)
	}

	// Create gRPC server instance
	var grpcServer = grpc.NewServer()

	// Register the server service on the gRPC server
	service.RegisterServiceServer(grpcServer, s)

	// Accept incoming connections on the listener
	s.logger.InfoPrintln("Server started.")
	err = grpcServer.Serve(listener)
	if err != nil {
		s.logger.ErrorLogger.Fatalf("Failed to start gRPC server. :: %v", err)
	}
}

func NewServer(logger *utils.Logger, address string, port string) *Server {
	return &Server{
		ipAddress: address+":"+port,
		logger: logger,
	}
}

