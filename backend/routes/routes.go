package routes

import (
	"compro-backend/db"
	"compro-backend/handlers"
	"compro-backend/middleware"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, database *db.DB, s3Client *s3.Client) {
	contentHandler := handlers.NewContentHandler(database)
	authHandler := handlers.NewAuthHandler(database)
	uploadHandler := handlers.NewUploadHandler(s3Client)
	newsHandler := handlers.NewNewsHandler(database)

	api := r.Group("/api")
	{
		// Public
		api.GET("/content", contentHandler.GetContent)
		api.GET("/news", newsHandler.GetPosts)
		api.GET("/news/:id", newsHandler.GetPost)

		// Auth
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/seed", authHandler.SeedAdmin)

		// Protected (any authenticated user)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.PUT("/content", contentHandler.UpdateContent)
			protected.POST("/upload", uploadHandler.UploadImage)

			// News management (admin + editor)
			news := protected.Group("")
			news.Use(middleware.RequireRole("admin", "editor"))
			{
				news.POST("/news", newsHandler.CreatePost)
				news.PUT("/news/:id", newsHandler.UpdatePost)
				news.DELETE("/news/:id", newsHandler.DeletePost)
			}

			// User management (admin only)
			users := protected.Group("")
			users.Use(middleware.RequireRole("admin"))
			{
				users.POST("/auth/users", authHandler.CreateUser)
				users.GET("/auth/users", authHandler.ListUsers)
				users.DELETE("/auth/users/:id", authHandler.DeleteUser)
			}
		}
	}
}
