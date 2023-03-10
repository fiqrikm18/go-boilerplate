package cmd

import (
	"github.com/fiqrikm18/go-boilerplate/internal/config"
	grpcRouter "github.com/fiqrikm18/go-boilerplate/internal/router/grpc"
	httpRouter "github.com/fiqrikm18/go-boilerplate/internal/router/http"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "",
	Short: "",
	Run:   runRootCommand,
}

func Execute() {
	err := rootCommand.Execute()
	if err != nil {
		panic(err)
	}
}

func runRootCommand(cmd *cobra.Command, args []string) {
	go func() {
		httpServer, err := config.NewHttpServer()
		if err != nil {
			panic(err)
		}

		httpRouter.RegisterRouter(httpServer.Srv)
		httpServer.Run()
	}()

	go func() {
		grpcServer, err := config.NewGrpcServer()
		if err != nil {
			panic(err)
		}

		grpcRouter.RegisterGRPCService(grpcServer.Srv)
		if err := grpcServer.Run(); err != nil {
			panic(err)
		}
	}()

	select {}
}
