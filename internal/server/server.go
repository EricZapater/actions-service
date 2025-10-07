package server

import (
	"actions-service/internal/setup"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(app *setup.App) {	
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

	serverHandlers := NewHandler(app)
	api := server.Group("/api")
	api.GET("/healthcheck", serverHandlers.HealthCheck)
	api.GET("/reload", serverHandlers.ReloadDTO)

	//api := server.Group("/api")
	//HealthCheck
	//api.GET("/healthcheck", controllers.HealthCheck.HealthCheck)
	//Status
	//api.POST("/status", controllers.Status.UpdateWorkcenterStatus)
	//Operator
	/*api.POST("/operator/clockin", controllers.Operator.ClockIn)
	api.POST("/operator/clockout", controllers.Operator.ClockOut)*/

	
	if err := server.Run(fmt.Sprintf(":%s", app.Cfg.ApiPort)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}