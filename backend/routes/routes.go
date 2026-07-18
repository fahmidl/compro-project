package routes

import (
	"compro-backend/handlers"
	"compro-backend/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(r *gin.Engine, db *mongo.Database, uploadDir string) {
	contentHandler := handlers.NewContentHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	uploadHandler := handlers.NewUploadHandler(uploadDir)

	api := r.Group("/api")
	{
		// Public
		api.GET("/content", contentHandler.GetContent)

		// Auth
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/seed", authHandler.SeedAdmin)

		// Protected
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.PUT("/content", contentHandler.UpdateContent)
			protected.POST("/upload", uploadHandler.UploadImage)
		}
	}
}
