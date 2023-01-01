package config

import (
	"fmt"

	"github.com/fiqrikm18/go-boilerplate/pkg/lib"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Srv     *gin.Engine
	AppConf *lib.ApplicationConfig
}

func NewHttpServer() (*HttpServer, error) {
	srv := gin.Default()
	conf, err := lib.LoadConfigFile()
	if err != nil {
		return nil, err
	}

	return &HttpServer{
		Srv:     srv,
		AppConf: conf,
	}, nil
}

func (h *HttpServer) Run() {
	addrPort := fmt.Sprintf(":%d", h.AppConf.HttpPort)

	if h.AppConf.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.ForceConsoleColor()
	h.Srv.Use(gin.Logger())
	h.Srv.Run(addrPort)
}
