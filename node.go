package main

import (
	"context"
	"google.golang.org/grpc"
	utils
)

func main() {

}

type ServiceClient interface {
	BroadcastRequest(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Request, error)
	RespondNode(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

func Broadcast() {
	var request = &{
		lamport:
	}
}

func BroadcastRequest(ctx context.Context, )
