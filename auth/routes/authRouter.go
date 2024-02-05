package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/hayk2377/distributed-ludo/auth/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", controller.Signup())
	incomingRoutes.POST("/login", controller.Login())
}
