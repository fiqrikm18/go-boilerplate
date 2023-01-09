package auth

import (
	"github.com/fiqrikm18/go-boilerplate/internal/model/dao"
	"github.com/fiqrikm18/go-boilerplate/pkg/lib"
	"net/http"
	"strings"
	"time"

	"github.com/fiqrikm18/go-boilerplate/internal/model/dto"
	"github.com/fiqrikm18/go-boilerplate/internal/repository"
	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	userRepository  *repository.UserRepository
	tokenRepository *repository.OAuthAccessTokenRepository
}

func NewAuthenticationController() *AuthenticationController {
	userRepository, err := repository.NewUserRepository()
	if err != nil {
		panic(err)
	}

	tokenRepository, err := repository.NewOAuthAccessTokenRepository()
	if err != nil {
		panic(err)
	}

	return &AuthenticationController{
		userRepository,
		tokenRepository,
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

	tokenUtils, err := lib.NewJWTToken()
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

	if err := controller.tokenRepository.RevokeByUserID(user.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	tokenDAO := dao.OauthAccessToken{
		AccessTokenUUID:     tokenData.AccessTokenUUID,
		RefreshTokenUUID:    tokenData.RefreshTokenUUID,
		AccessTokenExpDate:  time.Unix(tokenData.AccessTokenExpire, 0),
		RefreshTokenExpData: time.Unix(tokenData.RefreshTokenExpire, 0),
		Revoked:             false,
		UserId:              user.Id,
	}

	if err := controller.tokenRepository.Create(tokenDAO); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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

func (controller *AuthenticationController) Logout(c *gin.Context) {
	authToken := strings.Split(c.Request.Header["Authorization"][0], " ")
	if len(authToken) > 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid authentication token",
		})
		return
	}

	tokenUtils, err := lib.NewJWTToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := tokenUtils.ExtractToken(authToken[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := controller.userRepository.FindByUsernameOrEmail(token.Username, token.Username)
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

	if err := controller.tokenRepository.RevokeByUserID(user.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout success",
	})
}

func (controller *AuthenticationController) RefreshToken(c *gin.Context) {
	authToken := strings.Split(c.Request.Header["Authorization"][0], " ")
	if len(authToken) > 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid authentication token",
		})
		return
	}

	tokenUtils, err := lib.NewJWTToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := tokenUtils.ExtractToken(authToken[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := controller.userRepository.FindByUsernameOrEmail(token.Username, token.Username)
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

	tokenPayload := lib.TokenPayload{Username: user.Username}
	tokenData, err := tokenUtils.GenerateToken(tokenPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := controller.tokenRepository.RevokeByUserID(user.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	tokenDAO := dao.OauthAccessToken{
		AccessTokenUUID:     tokenData.AccessTokenUUID,
		RefreshTokenUUID:    tokenData.RefreshTokenUUID,
		AccessTokenExpDate:  time.Unix(tokenData.AccessTokenExpire, 0),
		RefreshTokenExpData: time.Unix(tokenData.RefreshTokenExpire, 0),
		Revoked:             false,
		UserId:              user.Id,
	}

	if err := controller.tokenRepository.Create(tokenDAO); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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
