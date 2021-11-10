package client

import (
	"google.golang.org/grpc"
	"mandatory-exercise-2/service"
	"mandatory-exercise-2/utils"
)

type Client struct {
	ipAddress string
	logger *utils.Logger
	Client service.ServiceClient
}

func (c *Client) ConnectClient() {
	c.logger.InfoPrintln("Connecting client...")
	// Create client connection to server
	conn, err := grpc.Dial(c.ipAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		c.logger.ErrorLogger.Fatalf("Could not connect to server. :: %v", err)
	}

	c.Client = service.NewServiceClient(conn)
	c.logger.InfoPrintln("Client connected.")
}

func NewClient(logger *utils.Logger, address string, port string)  *Client {
	return &Client{
		ipAddress: address+":"+port,
		logger: logger,
	}
}
