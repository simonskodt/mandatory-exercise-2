package main

import (
	"flag"
	"github.com/hashicorp/serf/serf"
	"mandatory-exercise-2/client"
	"mandatory-exercise-2/server"
	"mandatory-exercise-2/utils"
	"os"
	"strconv"
)

var Node *node
var logger *utils.Logger
var done = make(chan int)

type node struct {
	id int
	address string
	port string
	client *client.Client
	server *server.Server
	lamport *utils.Lamport
}

func main() {
	var id = flag.Int("id", 0, "The id of the node.") // Remember to input unique id
	var address = flag.String("address", "localhost", "The ip address of the node.")
	var port = flag.String("port", "8080", "The port of the node.")
	flag.Parse()

	// Create logger
	logger = utils.NewLogger("node_"+ strconv.Itoa(*id))

	// Create node and start node server
	logger.InfoPrintln("Creating node...")
	Node = newNode(*id, *address, *port)
	Node.init()
	logger.InfoPrintln("Node created.")

	//Setup cluster
	cluster, err := setupCluster(os.Getenv("ADVERTISE_ADDR"), os.Getenv("CLUSTER_ADDR"))
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed setting up cluster. :: %v", err)
	}
	defer func(cluster *serf.Serf) {
		err := cluster.Leave()
		if err != nil {
			logger.ErrorLogger.Fatalf("Fatal error. :: %v", err)
		}
	}(cluster)

	<- done
}

func (n *node) init() {
	go n.server.StartServer()
	n.client.ConnectClient()
}

func setupCluster(advertiseAddr string, clusterAddr string) (*serf.Serf, error){
	conf := serf.DefaultConfig()
	conf.Init()
	conf.MemberlistConfig.AdvertiseAddr = advertiseAddr

	cluster, err := serf.Create(conf)
	if err != nil {
		return nil, err
	}

	_, err = cluster.Join([]string{clusterAddr}, true)
	if err != nil {
		logger.ErrorLogger.Printf("Could not join cluster, starting own. :: %v\n", err)
	}

	return cluster, nil
}

func newNode(id int, address string, port string) *node {
	return &node{
		id: id,
		address: address,
		port: port,
		server: server.NewServer(logger, address, port),
		client: client.NewClient(logger, address, port),
		lamport: &utils.Lamport{T: 0},
	}
}