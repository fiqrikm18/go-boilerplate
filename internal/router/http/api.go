package http

import (
	AuthenticationController "github.com/fiqrikm18/go-boilerplate/internal/controller/auth"
	HomeController "github.com/fiqrikm18/go-boilerplate/internal/controller/home"
	"github.com/fiqrikm18/go-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(srv *gin.Engine) {
	srv.GET("", HomeController.IndexController)

	authController := AuthenticationController.NewAuthenticationController()
	srv.POST("/register", authController.RegisterController)
	srv.POST("/login", authController.LoginController)
	srv.GET("/logout", middleware.AuthMiddleware(), authController.Logout)
	srv.GET("/refresh-token", middleware.AuthMiddleware(), authController.RefreshToken)
}
