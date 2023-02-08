package grpc

import (
	"github.com/fiqrikm18/go-boilerplate/internal/model/protobuf"
	"golang.org/x/net/context"
)

type RpcGreetingServer struct {
	protobuf.UnimplementedGreeterServer
}

func (service *RpcGreetingServer) SayHello(ctx context.Context, req *protobuf.HelloRequest) (*protobuf.HelloReply, error) {
	return &protobuf.HelloReply{Message: "Hello " + req.Name}, nil
}

func (service *RpcGreetingServer) mustEmbedUnimplementedGreeterServer() {}
