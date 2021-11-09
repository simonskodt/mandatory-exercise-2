package main

import (
	"context"
	""
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
