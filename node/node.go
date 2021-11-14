package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"service"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"utils"
)

// DeleteLogsAfterExit determines whether to delete the log file for this node or not at program termination.
const DeleteLogsAfterExit = true

var node *Node
var done = make(chan int)

type Node struct {
	name    string	// id
	address string
	server  *server
	client  *client
	peers   map[string]service.ServiceClient
	logger  *utils.Logger
	lamport *utils.Lamport
	queue *utils.Queue
	state   int
	time    int
	service.UnimplementedServiceServer
}

func main() {
	var name = flag.String("name", "node0", "The unique name of the node.")
	var serverPort = flag.Int("sport", 8080, "The server port.")
	var ports = flag.String("ports", "", "The ports to the other servers.")
	var timer = flag.Int("time", 0, "The delay start time.")
	flag.Parse()

	node = NewNode(*name, "localhost", *serverPort, *timer)

	// Setup close handler for CTRL + C
	setupCloseHandler()

	otherPorts := strings.Split(*ports, ":")
	for _, port := range otherPorts {
		p, _ := strconv.Atoi(port)
		go node.registerPeers(p)
	}

	node.Init()

	<-done
	node.logger.WarningPrintln("NODE STOPPED.")
	os.Exit(0)
}

func (n *Node) Init() {
	// start internal server and connect internal client to server.
	n.start()
	n.runAlgorithm()
}

func (n *Node) GetName(context.Context, *service.Send) (*service.Info, error) {
	return &service.Info{Name: n.name}, nil
}

func (n *Node) registerPeers(port int) {
	conn, err := grpc.Dial("localhost"+":"+strconv.Itoa(port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		n.logger.ErrorFatalf("Could not create connection to other node. :: %v", err)
	}

	c := service.NewServiceClient(conn)
	info, _ := c.GetName(context.Background(), &service.Send{})
	n.peers[info.Name] = service.NewServiceClient(conn)
}

// start the internal server and connect the internal client to it.
func (n *Node) start() {
	n.logger.WarningPrintln("STARTING NODE...")

	go n.server.start(n.address, n.server.port)
	n.client.connect(n.address, n.server.port)

	n.logger.WarningPrintln("NODE STARTED.")
}

func (n *Node) runAlgorithm() {
	time.Sleep(time.Duration(n.time) * time.Second)
	n.state = utils.WANTED
	fmt.Printf("NODE %v entered WANTED\n", n.name)

	// Multicast
	replies, _ := n.BroadcastRequest(context.Background(), &service.Request{})

	if replies.Replies == int32(len(n.peers)) {
		n.state = utils.HELD
		time.Sleep(5 * time.Second)
		n.exit()
	} else {
		fmt.Println("NOT ENOUGH", replies.Replies, " LENGTH:", len(n.peers))
	}
}

func (n *Node) exit() {
	n.state = utils.RELEASED
	for n.queue.IsEmpty() {
		tuple := n.queue.Dequeue()
		n.peers[tuple.Name].ReplySender(context.Background(), &service.Request{})
	}
}

func (n *Node) BroadcastRequest(ctx context.Context, r *service.Request) (*service.Reply, error) {
	wait := sync.WaitGroup{}
	doneSending := make(chan int)
	var counter int32

	// INCREMENT LAMPORT
	n.lamport.Increment()

	for cname, c := range n.peers {
		wait.Add(1)
		go func() {
			defer wait.Done()
			fmt.Printf("NODE %v is broadcasting to %v\n", n.name, cname)
			_, err := c.Publish(ctx, &service.Request{Id: n.name, Lamport: n.lamport.T})
			fmt.Println(err)
			if err != nil {
				fmt.Println("Message received counter up!")
				counter++
			}
		}()
	}

	go func() {
		wait.Wait()
		close(doneSending)
	}()

	<- doneSending
	return &service.Reply{Replies: counter}, nil
}

func (n *Node) Publish(ctx context.Context, r *service.Request) (*service.Done, error) {
	n.receive(r.Lamport, r.Id)

	return &service.Done{}, nil
}

func (n *Node) ReplySender(ctx context.Context, r *service.Request) (*service.Done, error) {
	fmt.Printf("REPLY_SENDER (%v, Receive) NODE %v is replying %v\n", n.lamport.T, n.name, r.Id)
	return &service.Done{}, nil
}

func (n *Node) receive(lamport int32, name string) {
	if n.state == utils.HELD || (n.state == utils.WANTED && n.lamport.CompareLamportAndProcess(n.lamport.T, n.name, lamport, name)) {
		n.queue.Enqueue(lamport, name)
		n.lamport.MaxAndIncrement(lamport)
		fmt.Printf("IF (%v, Receive) NODE %v is enqueueing %v\n", n.lamport.T, n.name, name)
		// TODO: Observer for evt. fejl ved MaxAndIncrement
	} else {
		n.lamport.MaxAndIncrement(lamport) // Receive
		n.lamport.Increment() // Send
		n.ReplySender(context.Background(), &service.Request{Lamport: lamport, Id: name})
		fmt.Printf("ELSE (%v, Receive) NODE %v is replying %v\n", n.lamport.T, n.name, name)
	}
}

// setupCloseHandler sets a close handler for this program if it is interrupted
func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		node.server.stop()
		if DeleteLogsAfterExit {
			node.logger.DeleteLog()
		}
		close(done)
	}()
}

// NewNode creates a new node with the specified unique name and ip address.
func NewNode(name string, address string, serverPort int, time int) *Node {
	logger := utils.NewLogger(name)
	logger.WarningPrintf("CREATING NODE WITH ID '%v' AND IP ADDRESS '%v'", name, address)

	return &Node{
		name:    name,
		address: address,
		server:  newServer(serverPort, logger),
		client:  newClient(logger),
		peers:   make(map[string]service.ServiceClient),
		state:   utils.RELEASED,
		time:    time,
		logger:  logger,
		lamport: utils.NewLamport(),
		queue: utils.NewQueue(),
	}
}
