package auth

import (
	"github.com/fiqrikm18/go-boilerplate/pkg/lib"
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
	var request dto.RegisterRequest
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

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
	})
}

func (controller *AuthenticationController) LoginController(c *gin.Context) {
	var request dto.LoginRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := controller.userRepository.FindByUsernameOrEmail(request.Username, request.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if user.Email == "" || user.Username == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "invalid username or password",
		})
		return
	}

	tokenUtils, err := lib.NewToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	tokenPayload := lib.TokenPayload{Username: user.Username}
	tokenData, err := tokenUtils.GenerateToken(tokenPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := dto.LoginResponse{
		ExpiredIn:    tokenData.AccessTokenExpire,
		AccessToken:  tokenData.AccessToken,
		RefreshToken: tokenData.RefreshToken,
	}

	c.JSON(http.StatusOK, response)
}
