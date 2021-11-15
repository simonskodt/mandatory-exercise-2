package client

import (
	"google.golang.org/grpc"
	"service"
	"utils"
)

func NewClient(ipAddress string, logger *utils.Logger) service.ServiceClient {
	logger.InfoPrintf("Trying to connect to peer at %v.", ipAddress)

	// Create client connection to a server
	conn, err := grpc.Dial(ipAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.ErrorFatalf("Could not connect to peer. :: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logger.ErrorFatalf("Could not close connection of peer at %v. :: %v", ipAddress, err)
		}
	}(conn)

	logger.InfoPrintln("Successfully connected to peer.")

	return service.NewServiceClient(conn)
}
