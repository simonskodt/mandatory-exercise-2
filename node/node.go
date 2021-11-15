package main

import (
	"client"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"server"
	"service"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"utils"
)

// deleteLogsAfterExit determines whether to delete the log file for this node or not at program termination.
const deleteLogsAfterExit = true

// defaultAddress is the default address of all nodes.
const defaultAddress = "localhost"

// These constants represent the different states of a Node.
const (
	WANTED   int = 0
	HELD         = 1
	RELEASED     = 2
)

// The node of this process.
var node *Node
var done = make(chan int)

// A Node is a single process running on an ip address.
// It can communicate with other nodes.
type Node struct {
	name    string                           // name is the id of the Node.
	address string                           // address is the address of the Node.
	port    int                              // port is the port of the Node.
	server  *server.Server                   // server is the internal server.Server of the Node.
	peers   map[string]service.ServiceClient // peers is a map of all the other nodes in the cluster, mapping a Node name, to a service.ServiceClient.
	logger  *utils.Logger                    // logger is a log which logs specified
	lamport *utils.Lamport                   // lamport is a logical clock.
	queue   *utils.Queue                     // queue is the Node's FIFO queue of other nodes.
	state   int                              // The current state of the Node.
	time    int                              // Delay time before entering state WANTED.
	service.UnimplementedServiceServer
}

func main() {
	var name = flag.String("name", "node0", "The unique name of the node.")
	var serverPort = flag.Int("sport", 8080, "The server port.")
	var ports = flag.String("ports", "", "The ports to the other servers.")
	var timer = flag.Int("time", 0, "The delay start time.")
	flag.Parse()

	node = NewNode(*name, defaultAddress, *serverPort, *timer)

	// Setup close handler for CTRL + C
	setupCloseHandler()

	node.start(*ports)

	<-done
	node.server.Stop()
	node.logger.WarningPrintln("NODE STOPPED.")
	os.Exit(0)
}

// registerPeer connects to and register another Node on this node.
func (n *Node) registerPeer(port int) {
	peer := client.NewClient(createIpAddress(defaultAddress, port), n.logger)
	info, _ := peer.GetName(context.Background(), &service.InfoRequest{})
	n.peers[info.Name] = peer
}

// start the internal server and connect to other peers.
func (n *Node) start(ports string) {
	n.logger.WarningPrintln("STARTING NODE...")

	n.server.Start(n.getIpAddress(), n)

	n.logger.WarningPrintln("NODE STARTED.")

	otherPorts := strings.Split(ports, ":")
	for _, port := range otherPorts {
		p, _ := strconv.Atoi(port)
		go node.registerPeer(p)
	}

	// Wait before entering WANTED.
	time.Sleep(time.Duration(n.time) * time.Second)
	n.enter()
}

func (n *Node) enter() {
	n.state = WANTED
	n.logger.InfoPrintf("%v entered WANTED\n", n.name)

	// Multicast and wait for replies
	replies := n.multicast()

	if replies == len(n.peers) {
		n.state = HELD
		n.logger.InfoPrintf("%v entered HELD\n", n.name)
		time.Sleep(5 * time.Second)
		n.exit()
	} else {
		n.logger.WarningPrintf("%v did not get enough replies: %v/%v", replies, len(n.peers))
	}
}

// multicast sends a request to all peers of a Node and counts and returns the number of replies.
func (n *Node) multicast() int {
	n.logger.InfoPrintf("%v is now multicasting to peers.", n.name)

	wait := sync.WaitGroup{}
	doneSending := make(chan int)

	// Thread safe counter
	replies := utils.NewCounter()

	// Increment Lamport
	n.lamport.Increment()

	for clientName, serviceClient := range n.peers {
		wait.Add(1)

		go func(receiverName string, client service.ServiceClient) {
			defer wait.Done()
			n.logger.InfoPrintf("(%v, Send) %v is sending a request to %v.\n", n.lamport.Value(), n.name, receiverName)

			_, err := client.Publish(context.Background(), &service.Request{Lamport: n.lamport.Value(), Id: n.name})
			if err != nil {
				n.logger.ErrorPrintf("Error sending multicast to %v. :: %v\n", receiverName, err)
			} else {
				n.logger.InfoPrintf("Successfully send request to %v.\n", receiverName)
				replies.Increment()
			}

		}(clientName, serviceClient)
	}

	go func() {
		wait.Wait()
		close(doneSending)
	}()

	<-doneSending
	n.logger.InfoPrintf("%v done multicasting. Replies: %v", n.name, replies.Value())

	return replies.Value()
}

// receive a service.Request from a Node and either reply back to the Node or enqueue it in the queue.
func (n *Node) receive(lamport int32, name string) error {
	n.logger.InfoPrintf("%v received request from %v.\n", n.name, name)

	if n.state == HELD || (n.state == WANTED && n.lamport.CompareLamportAndProcess(n.name, lamport, name)) {
		n.queue.Enqueue(lamport, name)
		n.lamport.MaxAndIncrement(lamport)
		// TODO: Observer for evt. fejl ved MaxAndIncrement
		return fmt.Errorf("(%v, Receive) %v is enqueing %v\n", n.lamport.Value(), n.name, name)
	}

	n.lamport.MaxAndIncrement(lamport) // Receive
	n.lamport.Increment()              // Send
	n.logger.InfoPrintf("(%v, Receive) %v is replying %v -> GO AHEAD!\n", n.lamport.Value(), n.name, name)
	return nil
}

func (n *Node) exit() {
	n.state = RELEASED
	for !n.queue.IsEmpty() {
		_, name := n.queue.Dequeue()
		n.peers[name].ReplySender(context.Background(), &service.Request{})
	}
}

// Publish receives requests from another node.
func (n *Node) Publish(_ context.Context, r *service.Request) (*service.Reply, error) {
	reply := n.receive(r.Lamport, r.Id)
	return &service.Reply{}, reply
}

func (n *Node) ReplySender(_ context.Context, r *service.Request) (*service.Reply, error) {
	fmt.Printf("REPLY_SENDER (%v, Receive) NODE %v is replying %v\n", n.lamport.Value(), n.name, r.Id)
	return &service.Reply{}, nil
}

func (n *Node) GetName(context.Context, *service.InfoRequest) (*service.InfoReply, error) {
	return &service.InfoReply{Name: n.name}, nil
}

// setupCloseHandler sets a close handler for this program if it is interrupted
func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		if deleteLogsAfterExit {
			node.logger.DeleteLog()
		}
		<-c
		close(done)
	}()
}

func (n *Node) getIpAddress() string {
	return n.address + ":" + strconv.Itoa(n.port)
}

func createIpAddress(address string, port int) string {
	return address + ":" + strconv.Itoa(port)
}

// NewNode creates a new node with the specified unique name and ip address.
func NewNode(name string, address string, serverPort int, time int) *Node {
	logger := utils.NewLogger(name)
	logger.WarningPrintf("CREATING NODE WITH ID '%v' AND IP ADDRESS '%v'", name, address)

	return &Node{
		name:    name,
		address: address,
		port:    serverPort,
		server:  server.NewServer(logger),
		lamport: utils.NewLamport(),
		queue:   utils.NewQueue(),
		peers:   make(map[string]service.ServiceClient),
		state:   RELEASED,
		time:    time,
		logger:  logger,
	}
}
