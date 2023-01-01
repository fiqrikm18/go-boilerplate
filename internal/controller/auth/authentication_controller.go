package auth

import (
	"net/http"

	"github.com/fiqrikm18/go-boilerplate/internal/model/dto"
	"github.com/fiqrikm18/go-boilerplate/internal/repository"
	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	userRepository *repository.UserRepository
}

func NewAuthenticationController() *AuthenticationController {
	userRepository, err := repository.NewUserRepository()
	if err != nil {
		panic(err)
	}

	return &AuthenticationController{
		userRepository,
	}
}

func (controller *AuthenticationController) RegisterController(c *gin.Context) {
	var request dto.UserRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"messagee": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"asd": request,
	})
}
