package server

import (
	"actions-service/internal/config"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {	
	fmt.Println("Starting server...")
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		/*AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},*/
		MaxAge: 12 * time.Hour,
	}))



	//api := server.Group("/api")
	//HealthCheck
	//api.GET("/healthcheck", controllers.HealthCheck.HealthCheck)
	//Status
	//api.POST("/status", controllers.Status.UpdateWorkcenterStatus)
	//Operator
	/*api.POST("/operator/clockin", controllers.Operator.ClockIn)
	api.POST("/operator/clockout", controllers.Operator.ClockOut)*/

	
	if err := server.Run(fmt.Sprintf(":%s", cfg.ApiPort)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}