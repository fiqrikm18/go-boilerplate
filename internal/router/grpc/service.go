package grpc

import (
	grpcService "github.com/fiqrikm18/go-boilerplate/internal/grpc"
	"github.com/fiqrikm18/go-boilerplate/internal/model/protobuf"
	"google.golang.org/grpc"
)

func RegisterGRPCService(srv *grpc.Server) {
	protobuf.RegisterGreeterServer(srv, &grpcService.RpcGreetingServer{})
}
