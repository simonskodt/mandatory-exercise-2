package main

import (
	"context"
	"flag"
	"mandatory-exercise-2/client"
	"mandatory-exercise-2/server"
	"mandatory-exercise-2/service"
	"mandatory-exercise-2/utils"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// deleteLogsAfterExit determines whether to delete the log file for this node or not at program termination.
const deleteLogsAfterExit = true

// defaultAddress is the default address of all nodes.
const defaultAddress = "localhost"

// Constants which represent the different states of a node.
const (
	WANTED   int = 0
	HELD         = 1
	RELEASED     = 2
)

var done = make(chan int)
var wait sync.WaitGroup

// A node is a single process running on an ip address.
// It can communicate with other nodes.
type node struct {
	name       string                           // name is the id of the node.
	ipAddress  *net.TCPAddr                     // ipAddress is the full ip address of the node.
	state      int                              // The current state of the node.
	server     *server.Server                   // server is the internal server.Server of the node.
	logger     *utils.Logger                    // logger is a log which logs specified
	lamport    *utils.Lamport                   // lamport is a logical clock.
	queue      *utils.Queue                     // queue is the node's FIFO queue of other nodes.
	repCounter *utils.Counter                   // repCounter counts the number of replies from other peers.
	peers      map[string]service.ServiceClient // peers is a map of all the other nodes in the cluster, mapping a node name, to a service.ServiceClient.
	delay      time.Duration                    // delay time before entering state WANTED.
	service.UnimplementedServiceServer
}

func main() {
	var name = flag.String("name", "node0", "The unique name of the node.")
	var address = flag.String("address", defaultAddress, "The address of the node.")
	var serverPort = flag.Int("sport", 8080, "The server port.")
	var ports = flag.String("ports", "", "The ports to the other nodes.")
	var delay = flag.Duration("delay", 0, "The delay time in seconds before the node enters enter().") //flag.Int("delay", 0, "The delay start time.")

	// Setup close handler for CTRL + C
	setupCloseHandler()

	//Create and start the node.
	n := newNode(*name, *address, *serverPort, *delay)
	n.start(*ports)

	<-done
	n.stop()
	if deleteLogsAfterExit {
		n.logger.DeleteLog()
	}
	os.Exit(0)
}

// registerPeer connects to and registers another node on this node at the specified port.
func (n *node) registerPeer(port int) {
	peer := client.NewClient(createIpAddress(defaultAddress, port), n.logger)
	info, err := peer.GetName(context.Background(), &service.NameRequest{Name: n.name})
	if err != nil {
		n.logger.ErrorFatalf("Could not fetch name of peer. :: %v", err)
	}
	n.peers[info.Name] = peer
}

// start the node and connect to other peers (nodes).
func (n *node) start(ports string) {
	n.logger.WarningPrintln("STARTING NODE...")
	n.server.Start(n.ipAddress.String(), n)
	n.logger.WarningPrintln("NODE STARTED.")

	peerPorts := strings.Split(ports, ":")
	for _, port := range peerPorts {
		p, _ := strconv.Atoi(port)
		go n.registerPeer(p)
	}

	// Wait before entering WANTED.
	time.Sleep(n.delay * time.Second)
	n.enter()
}

// stop shutdowns the node.
func (n *node) stop() {
	n.logger.WarningPrintln("STOPPING NODE...")
	n.server.Stop()
	n.logger.WarningPrintln("NODE STOPPED.")
}

// enter makes the node enter WANTED.
// It multicasts a request to all peers and waits for replies.
// N-1 replies must be received before entering HELD.
func (n *node) enter() {
	n.state = WANTED
	n.logger.InfoPrintf("%v entered WANTED\n", n.name)

	// Multicast to all peers
	replies := n.multicast()

	if replies == len(n.peers) {
		n.state = HELD
		n.logger.InfoPrintf("%v entered HELD\n", n.name)
		time.Sleep(5 * time.Second)
		n.exit()
	} else {
		n.logger.WarningPrintf("%v did not get enough repCounter: %v/%v", n.name, replies, len(n.peers))
	}
}

// multicast sends a request to all peers of a node and counts and returns the number of peers responding.
func (n *node) multicast() int {
	n.lamport.Increment()
	n.logger.InfoPrintf("(%v, Send) %v is now multicasting to peers.\n", n.lamport.Value(), n.name)

	wait = sync.WaitGroup{}
	doneSending := make(chan int)

	// Reset counter at each multicast call.
	n.repCounter.Reset()

	for clientName, serviceClient := range n.peers {
		wait.Add(1)

		go func(receiverName string, client service.ServiceClient) {
			n.logger.InfoPrintf("(%v, Send) %v is sending a request to %v.\n", n.lamport.Value(), n.name, receiverName)

			reply, err := client.Publish(context.Background(), &service.Request{Lamport: n.lamport.Value(), Name: n.name})
			if err != nil {
				n.logger.ErrorPrintf("Error sending request to %v. :: %v\n", receiverName, err)
			}

			if reply.Ack {
				n.replyReceived()
			}
		}(clientName, serviceClient)
	}

	go func() {
		wait.Wait()
		close(doneSending)
	}()

	<-doneSending
	n.logger.InfoPrintf("(%v) %v done multicasting. Replies: %v/%v\n", n.lamport.Value(), n.name, n.repCounter.Value(), len(n.peers))

	return n.repCounter.Value()
}

// replyReceived is called when a node receives a reply.
// It decrements wait and increments repCounter.
func (n *node) replyReceived() {
	n.repCounter.Increment()
	wait.Done()
}

// receive a service.Request from a node and either reply back to the node or enqueue it in the queue.
func (n *node) receive(lamport int32, name string) bool {
	n.logger.InfoPrintf("%v received request from %v.\n", n.name, name)

	if n.state == HELD || (n.state == WANTED && n.lamport.CompareLamportAndProcess(n.name, lamport, name)) {
		n.queue.Enqueue(lamport, name)
		n.logger.InfoPrintf("%v is enqueued %v\n", n.name, name)
		return false
	}

	n.logger.InfoPrintf("(%v, Receive) %v is replying %v -> GO AHEAD!\n", n.lamport.Value(), n.name, name)
	n.lamport.MaxAndIncrement(lamport) // Receive
	n.lamport.Increment()              // Send reply back
	return true
}

// exit releases the CS and sends a reply to all peers.
func (n *node) exit() {
	n.state = RELEASED
	n.logger.InfoPrintf("%v entered RELEASED\n", n.name)

	for !n.queue.IsEmpty() {
		_, name := n.queue.Dequeue()
		n.logger.InfoPrintf("%v dequeued %v\n", n.name, name)

		_, err := n.peers[name].ReplySender(context.Background(), &service.Request{Name: n.name, Lamport: n.lamport.Value()})
		if err != nil {
			n.logger.ErrorPrintf("Could not dequeue %v. :: %v\n", name, err)
		}
	}
}

// Publish receives requests from another node.
func (n *node) Publish(_ context.Context, r *service.Request) (*service.Reply, error) {
	reply := n.receive(r.Lamport, r.Name)
	return &service.Reply{Ack: reply}, nil
}

// ReplySender sends a reply from a node which is received at the peer.
func (n *node) ReplySender(_ context.Context, r *service.Request) (*service.Reply, error) {
	n.logger.InfoPrintf("(%v, Send) %v is replying %v.\n", n.lamport.Value(), r.Name, n.name)
	n.replyReceived()
	return &service.Reply{}, nil
}

// GetName returns an info struct to the caller.
func (n *node) GetName(_ context.Context, nq *service.NameRequest) (*service.NameReply, error) {
	n.logger.InfoPrintf("%v is requesting the name of %v.\n", nq.Name, n.name)
	n.logger.InfoPrintf("Sending back %v.\n", n.name)
	return &service.NameReply{Name: n.name}, nil
}

// createIpAddress converts an address string and a port (integer) to a string.
func createIpAddress(address string, port int) string {
	return address + ":" + strconv.Itoa(port)
}

// setupCloseHandler sets a close handler for this program if it is interrupted.
func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(done)
	}()
}

// newNode creates a new node with the specified unique name and ip address.
func newNode(name string, address string, serverPort int, delay time.Duration) *node {
	logger := utils.NewLogger(name)
	logger.WarningPrintf("CREATING NODE WITH ID '%v' AND IP ADDRESS '%v'", name, address)

	ipAddress, err := net.ResolveTCPAddr("tcp", createIpAddress(address, serverPort))
	if err != nil {
		logger.ErrorFatalf("Error resolving tcp address %v:%v. :: %v", address, serverPort, err)
	}

	return &node{
		name:       name,
		ipAddress:  ipAddress,
		state:      RELEASED,
		server:     server.NewServer(logger),
		lamport:    utils.NewLamport(),
		queue:      utils.NewQueue(),
		repCounter: utils.NewCounter(),
		peers:      make(map[string]service.ServiceClient),
		delay:      delay,
		logger:     logger,
	}
}
