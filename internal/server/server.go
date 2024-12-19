package server

import (
	"net/http"

	"github.com/did-server/docs"
	"github.com/did-server/internal/server/router"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewServer(chain, port string) *http.Server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(router.Cors())
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to DID Server",
		})
	})

	router.NewRouter(chain, r)

	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
}
