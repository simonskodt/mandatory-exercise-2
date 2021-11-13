package main

import (
	"google.golang.org/grpc"
	"service"
	"strconv"
	"utils"
)

type client struct {
	client service.ServiceClient
	logger *utils.Logger
}

func (c *client) connect(address string, port int) {
	c.logger.InfoPrintln("Connecting client...")
	// Create client connection to server
	conn, err := grpc.Dial(address + ":" + strconv.Itoa(port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		c.logger.ErrorLogger.Fatalf("Could not connect to server. :: %v", err)
	}
	defer conn.Close()

	c.client = service.NewServiceClient(conn)
	c.logger.InfoPrintln("Client connected.")
}

func newClient(logger *utils.Logger) *client {
	return &client{
		logger: logger,
	}
}
