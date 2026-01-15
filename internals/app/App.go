package app

import (
	"log/slog"
	"truthly/internals/routes"
	"truthly/internals/util"
	"truthly/internals/util/tail38"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start(logger *slog.Logger) {
	db := util.InitDb()
	logger.Info("DB Connected")

	client := tail38.InitGeoClient(
		"127.0.0.1",
		9851,
		logger,
	)

	sqlDb, err := db.DB()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer sqlDb.Close() // Close db connection when the app is shut down

	router := gin.Default()

	// ADD CORS Policy
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000", "http://127.0.0.1/*",
		},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Origin"},
		AllowCredentials: true,
	}))

	// health route if my app is runnig
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "app is running",
		})
	})

	// versioned API
	api := router.Group("/api/v1")

	// register all routes
	routes.RegisterAll(api, db, logger, client)

	logger.Info("Server started at :8181")
	router.Run(":8181")
}
