package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/hayk2377/distributed-ludo/auth/controllers"
	"github.com/hayk2377/distributed-ludo/auth/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}