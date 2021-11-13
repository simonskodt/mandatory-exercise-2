package main

import (
	"context"
	"flag"
	"service"
	"utils"
)

var n *node
var done chan int
var logger *utils.Logger
var lamport *utils.Lamport

type node struct {
	name    string
	address string
	server  *server
	client *client
	peers map[string]service.ServiceClient
	service.UnimplementedServiceServer
}

func main() {
	var name = flag.String("name", "node0", "The unique name of the node.")
	var address = flag.String("address", "localhost:8080", "The ip address of the node. Usage: 'ip:port'")
	flag.Parse()

	logger = utils.NewLogger(*name)
	lamport = utils.NewLamport()
	n = NewNode(*name, *address)
	done = make(chan int)

	// Start internal server and connect internal client
	n.start()

	<- done
}

func (n *node) start() {
	logger.WarningPrintln("STARTING NODE...")

	// Register the gRPC server on the node service
	service.RegisterServiceServer(n.server.grpcServer, n)

	// Start the internal server
	go n.server.start(n.address)

	// Connect internal client to internal server
	n.client.connect(n.address)

	logger.WarningPrintln("NODE STARTED.")
}

func (n *node) registerService() {

}

func (n *node) Broadcast() {

}

func (n *node) BroadcastRequest(context.Context, *service.Request) (*service.Request, error)  {
	return nil, nil
}

func NewNode(name string, address string) *node {
	logger.WarningPrintf("CREATING NODE WITH ID '%v' AND IP ADDRESS '%v'", name, address)
	return &node{
		name:    name,
		address: address,
		server:  newServer(logger),
		client:  newClient(logger),
		peers:   make(map[string]service.ServiceClient),
	}
}
