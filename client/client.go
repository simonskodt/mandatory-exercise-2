package client

import (
	"google.golang.org/grpc"
	"mandatory-exercise-2/service"
	"mandatory-exercise-2/utils"
)

func NewClient(ipAddress string, logger *utils.Logger) service.ServiceClient {
	logger.InfoPrintf("Trying to connect to peer at %v.\n", ipAddress)

	// Create client connection to a server
	conn, err := grpc.Dial(ipAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		defer conn.Close()
		logger.ErrorFatalf("Could not connect to peer at %v. :: %v", ipAddress, err)
	}

	logger.InfoPrintf("Successfully connected to peer at %v.\n", ipAddress)

	return service.NewServiceClient(conn)
}
