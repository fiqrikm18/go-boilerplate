package config

import (
	"fmt"
	"net"
	"strconv"

	"github.com/fiqrikm18/go-boilerplate/pkg/lib"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Srv     *grpc.Server
	AppConf *lib.ApplicationConfig
}

func NewGrpcServer() (*GRPCServer, error) {
	conf, err := lib.LoadConfigFile()
	if err != nil {
		return nil, err
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	return &GRPCServer{
		Srv:     grpcServer,
		AppConf: conf,
	}, nil
}

func (srv *GRPCServer) Run() error {
	fmt.Println("GRPC Server Running on :" + strconv.Itoa(srv.AppConf.GrpcPort))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.AppConf.GrpcPort))
	if err != nil {
		return err
	}

	return srv.Srv.Serve(listener)
}
