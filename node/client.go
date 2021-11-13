package main

import (
	"google.golang.org/grpc"
	"service"
	"utils"
)

type client struct {
	logger *utils.Logger
	client service.ServiceClient
}

func (c *client) connect(ipAddress string) {
	c.logger.InfoPrintln("Connecting client...")
	// Create client connection to server
	conn, err := grpc.Dial(ipAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		c.logger.ErrorLogger.Fatalf("Could not connect to server. :: %v", err)
	}
	defer conn.Close()

	c.client = service.NewServiceClient(conn)
	c.logger.InfoPrintln("Client connected.")
}

func newClient(logger *utils.Logger)  *client {
	return &client{
		logger: logger,
	}
}
