package main

import (
	"context"
	"log"
	"os"

	"compro-backend/db"
	"compro-backend/routes"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	ginEngine *gin.Engine
	adapter   *ginadapter.GinLambda
)

func init() {
	// Load .env in local dev; Lambda provides env vars automatically
	godotenv.Load()

	// Load AWS SDK config (uses default credentials chain)
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("failed to load AWS config:", err)
	}

	// Initialize DynamoDB client
	dynamoClient := dynamodb.NewFromConfig(cfg)
	database := db.NewDB(dynamoClient)

	// Initialize S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Build the Gin engine
	r := gin.Default()

	// CORS — use CLOUDFRONT_DOMAIN in production, wildcard in local dev
	allowedOrigin := os.Getenv("CLOUDFRONT_DOMAIN")
	if allowedOrigin != "" {
		allowedOrigin = "https://" + allowedOrigin
	} else {
		allowedOrigin = "*"
	}
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup API routes
	routes.SetupRoutes(r, database, s3Client)

	ginEngine = r
	adapter = ginadapter.New(ginEngine)
}

// Lambda handler — receives API Gateway proxy requests and delegates to Gin.
func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.ProxyWithContext(ctx, req)
}

func main() {
	if os.Getenv("AWS_LAMBDA") != "" || os.Getenv("AWS_EXECUTION_ENV") != "" || os.Getenv("_HANDLER") != "" {
		// Running inside AWS Lambda
		log.Println("Starting Lambda handler")
		lambda.Start(lambdaHandler)
	} else {
		// Local development
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Server running on port %s", port)
		ginEngine.Run(":" + port)
	}
}
