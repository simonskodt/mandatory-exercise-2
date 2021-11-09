package main

import (
	"context"
	"mandatory-exercise-2/service"
	"sync"
)

type server struct {
	streams *[]service.Service_CreateStreamServer
	service.UnimplementedServiceServer
}

func BroadcastRequest(ctx context.Context, r *service.Request) (*service.Request, error) {

}

func RespondNode(ctx context.Context, r *service.Request) (*service.Response, error) {

}

func CreateStream(r *service.Request, stream service.Service_CreateStreamServer) error {



	return nil
}

func (s *server) Broadcast(_ context.Context, msg *service.Request) (*service.Request, error) {
	wait := sync.WaitGroup{}
	doneSending := make(chan int)

	// Send message to each client connection
	for _, conn := range s.connections {
		wait.Add(1)

		go func(msg *pb.Message, conn *connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				logger.InfoPrintf("(%v, Send) Broadcasting message '%v' from '%v' to '%v'", msg.Lamport, msg.Content, msg.User.Name, conn.user.Name)

				if err != nil {
					logger.ErrorPrintf("Error with stream %v. :: %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(doneSending)
	}()

	<-doneSending
	return &pb.Done{}, nil
}
