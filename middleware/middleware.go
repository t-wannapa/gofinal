package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "November 10, 2009" {
		c.JSON(http.StatusUnauthorized, "unauthorized!!")
		c.Abort()
		return
	}
	c.Next()
}
