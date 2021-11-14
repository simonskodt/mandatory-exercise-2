package main

import (
	"context"
	"os"
	"os/signal"
	"service"
	"syscall"
	"utils"
)

// DeleteLogsAfterExit determines whether to delete the log file for this node or not at program termination.
const DeleteLogsAfterExit = true
const address = "127.0.0.1"

var n *node
var done = make(chan int)
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
	name := "node0"
	address := address
	port := 8080

	// Setup close handler for CTRL + C
	setupCloseHandler()

	// Create structs
	logger = utils.NewLogger(name)
	lamport = utils.NewLamport()
	n = newNode(name, address, port)

	// Start Serf, internal server and connect internal client to server.
	n.start()

	<- done
	logger.WarningPrintln("NODE STOPPED.")
	os.Exit(0)
}

// start the internal server and connect the internal client to it. Also initiates the service.SerfService for automatic service registration.
func (n *node) start() {
	logger.WarningPrintln("STARTING NODE...")

	go n.server.start(n.address, n.server.port)
	n.client.connect(n.address, n.server.port)

	logger.WarningPrintln("NODE STARTED.")
}

func (n *node) BroadcastRequest(context.Context, *service.Request) (*service.Request, error)  {
	return nil, nil
}

// setupCloseHandler sets a close handler for this program if it is interrupted
func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<- c
		n.server.stop()
		if DeleteLogsAfterExit {
			logger.DeleteLog()
		}
		close(done)
	}()
}

// newNode creates a new node with the specified unique name and ip address.
func newNode(name string, address string, serverPort int) *node {
	logger.WarningPrintf("CREATING NODE WITH ID '%v' AND IP ADDRESS '%v'", name, address)
	return &node{
		name:    name,
		address: address,
		server:  newServer(serverPort, logger),
		client:  newClient(logger),
		peers:   make(map[string]service.ServiceClient),
	}
}
