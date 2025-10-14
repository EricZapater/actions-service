package server

import (
	"actions-service/internal/setup"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// Importar els packages de Swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Importar els docs generats (ajusta la ruta segons el teu projecte)
	_ "actions-service/docs"
)

// @title Actions Service API
// @version 1.0
// @description API per gestionar operadors i workcenters
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https
func Run(app *setup.App) {	
	fmt.Println("Starting server...")
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))
	
	// Ruta per Swagger UI
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	serverHandlers := NewHandler(app)
	api := server.Group("/api")
	api.GET("/healthcheck", serverHandlers.HealthCheck)
	api.GET("/reload", serverHandlers.ReloadDTO)
	api.POST("/operator/clockin", app.Handlers.OperatorHandler.ClockIn)
	api.POST("/operator/clockout", app.Handlers.OperatorHandler.ClockOut)
	ws := server.Group("/ws")
	ws.GET("/general", serverHandlers.WSGeneral)
	ws.GET("/workcenter/:id", serverHandlers.WSWorkcenter)
	
	if err := server.Run(fmt.Sprintf(":%s", app.Cfg.ApiPort)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}