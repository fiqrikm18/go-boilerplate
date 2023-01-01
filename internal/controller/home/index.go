package home

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "welcome to this simple golang boilerpalte",
	})
}
