package main

import (
	"flag"
	"github.com/hashicorp/serf/serf"
	"google.golang.org/grpc"
	"mandatory-exercise-2/service"
	"mandatory-exercise-2/utils"
	"net"
)

var Node *node

type node struct {
	id int
	address string
	port string
	client service.ServiceClient
	server service.ServiceServer
	lamport *utils.Lamport
	logger *utils.Logger
}

func main() {
	var id = flag.Int("id", 0, "The id of the node.") // Remember to input unique id
	var address = flag.String("address", "localhost", "The ip address of the node.")
	var port = flag.String("port", "8080", "The port of the node.")
	flag.Parse()


	Node = newNode(*id, *address, *port)
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logger.ErrorLogger.Fatalf("Could not connect to server. :: %v", err)
	}
	var _ = service.NewServiceClient(conn)
}

func (n *node) broadcast() {

}

func (n *node) respondClient()  {
	
}

func newNode(id int, address string, port string) *node {
	return &node{
		id: id,
		address: address,
		port: port,
		server: newServer(port),
		client: newClient(address, port),
		lamport: &utils.Lamport{T: 0},
		logger: utils.NewLogger("Node"),
	}
}

func newClient(address string, port string) service.ServiceClient {
	conn, err := grpc.Dial(address+":"+port)
	if err != nil {
		Node.logger.ErrorLogger.Fatalf("Could not connect to server. :: %v", err)
	}
	return service.NewServiceClient(conn)
}

func newServer(port string) service.ServiceServer {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		Node.logger.ErrorLogger.Fatalf("Could not listen at port %v. :: %v", port, err)
	}

	var grpcServer = grpc.NewServer()
	service.RegisterServiceServer(grpcServer, )
}