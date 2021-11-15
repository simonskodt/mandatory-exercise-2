package server

import (
	"google.golang.org/grpc"
	"net"
	"service"
	"utils"
)

// A Server is an gRPC server.
// To start the Server, call the Start function.
// To stop the Server, call the Stop function.
type Server struct {
	grpcServer *grpc.Server  // grpcServer is the grpcServer running in the Server.
	logger     *utils.Logger // logger is log, used to log all activities of the Server.
}

// Start the Server for a node on the specified ip address on a new go routine.
func (s *Server) Start(ipAddress string, node service.ServiceServer) {
	go func() {
		s.logger.InfoPrintln("Starting server...")
		// start listener
		listener, err := net.Listen("tcp", ipAddress)
		if err != nil {
			s.logger.ErrorFatalf("Could not listen at ip address %v. :: %v", ipAddress, err)
		}

		// Register the gRPC server on the node service
		service.RegisterServiceServer(s.grpcServer, node)

		// Accept incoming connections on the listener
		s.logger.InfoPrintln("server started.")
		err = s.grpcServer.Serve(listener)
		if err != nil {
			s.logger.ErrorFatalf("Failed to start gRPC server. :: %v", err)
		}
	}()
}

// Stop immediately stops the Server.
func (s *Server) Stop() {
	s.logger.WarningPrintln("Stopping server...")
	s.grpcServer.Stop()
	s.logger.WarningPrintln("Server stopped.")
}

// NewServer creates and returns a new Server which will be run on a specific ip address.
func NewServer(logger *utils.Logger) *Server {
	return &Server{
		grpcServer: grpc.NewServer(),
		logger:     logger,
	}
}
