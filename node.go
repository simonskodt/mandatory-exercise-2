package main

import (
	"google.golang.org/grpc"
	"mandatory-exercise-2/utils"
	"mandatory-exercise-2/service"
)

var lamport = &utils.Lamport{}
var logger = utils.NewLogger("Node")

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logger.ErrorLogger.Fatalf("Could not connect to server. :: %v", err)
	}
	var _ = service.NewServiceClient(conn)
}

func broadcast() {
	
}

func respondClient()  {
	
}