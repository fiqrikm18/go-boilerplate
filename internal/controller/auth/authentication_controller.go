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
			"message": err.Error(),
		})
		return
	}

	userCheck, err := controller.userRepository.FindByUsernameOrEmail(request.Username, request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if userCheck.Email != "" || userCheck.Username != "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "username or email already registered",
		})
		return
	}

	err = controller.userRepository.Create(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.UserResponse{
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
	})
}
