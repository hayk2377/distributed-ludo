package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	helper "github.com/hayk2377/distributed-ludo/auth/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")

		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		if strings.HasPrefix(clientToken, "Bearer ") {
			clientToken = strings.TrimPrefix(clientToken, "Bearer ")
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("id", claims.ID)
		c.Next()
	}
}
