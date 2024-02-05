package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	routes "github.com/hayk2377/distributed-ludo/auth/routes"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)

func main() {
	//Load env variable
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//get port
	port := os.Getenv("PORT")

	//if no port, pick 8000 by default
	if port == "" {
		port = "8000"
	}

	//get a new router, log every event
	router := gin.New()

  config := cors.DefaultConfig()

  config.AllowAllOrigins = true
  config.AllowHeaders = []string{"Authorization", "Content-Type"}

  router.Use(cors.New(config))

	router.Use(cors.Default())
	router.Use(gin.Logger())

	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
