package routes

import (
	"command-dispatcher/internal/core/interceptors"
	"command-dispatcher/internal/routes/command"
	"command-dispatcher/internal/routes/users"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

// Embed the web directory
//
// //go:embed web
// var static embed.FS

func Init() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	r := gin.New()
	// CORS configuration to allow all origins and expose all headers
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware registration
	r.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())
	r.Use(interceptors.JsonApiInterceptor())

	// Serving API
	api := r.Group("/api")

	// Routes registration
	users.Register(api)
	command.Register(api)

	// Start the Server
	log.Printf("Server is running on port: %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
