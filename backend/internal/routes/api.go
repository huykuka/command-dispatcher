package routes

import (
	"command-dispatcher/internal/core/interceptors"
	"command-dispatcher/internal/core/middlewares"
	"command-dispatcher/internal/routes/auth"
	"command-dispatcher/internal/routes/users"
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

// Embed the web directory
//
//go:embed web
var static embed.FS

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
	auth.Register(api)

	// Create a sub-filesystem for the "web" directory
	contentStatic, err := fs.Sub(static, "web")
	if err != nil {
		log.Fatal(err)
	}

	r.Use(middlewares.NoJsonAPI(), func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api/") {
			filePath := filepath.Join(".", c.Request.URL.Path)
			_, err := contentStatic.Open(filePath)
			c.Writer.Header().Del("Content-Length")
			if os.IsNotExist(err) {
				c.FileFromFS("index.html", http.FS(contentStatic))
			} else {

				http.StripPrefix("/", http.FileServer(http.FS(contentStatic))).ServeHTTP(c.Writer, c.Request)
			}
			c.Abort()
		}
	})

	// Start the Server
	log.Printf("Server is running on port: %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
