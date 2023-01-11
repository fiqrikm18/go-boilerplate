package middleware

import (
	"github.com/fiqrikm18/go-boilerplate/internal/model/dao"
	"github.com/fiqrikm18/go-boilerplate/internal/repository"
	"github.com/fiqrikm18/go-boilerplate/pkg/lib"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var (
	tokenRepository *repository.OAuthAccessTokenRepository
	tokenUtil       lib.JwtToken
	err             error
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenRepository, err = repository.NewOAuthAccessTokenRepository()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			context.Abort()
			return
		}

		if !strings.Contains(context.Request.Header["Authorization"][0], "Bearer") {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
			context.Abort()
			return
		}

		tokenUtil, err = lib.NewJWTToken()
		authToken := strings.Split(context.Request.Header["Authorization"][0], " ")
		if len(authToken) < 2 {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
			context.Abort()
			return
		}

		tokenClaims, err := tokenUtil.ExtractToken(authToken[0])
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			context.Abort()
			return
		}

		var tokenData dao.OauthAccessToken
		tx := tokenRepository.DbConn.DB.Find(&tokenData).Where("access_token_uuid=?", tokenClaims.TokenUUID)
		if err := tx.Error; err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			context.Abort()
			return
		}

		if tokenData.Revoked {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": "token not active",
			})
			context.Abort()
			return
		}

		if time.Now().After(tokenData.AccessTokenExpDate) {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": "token expired",
			})
			context.Abort()
			return
		}

		context.Next()
	}
}
