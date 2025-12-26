package main

import (
	"context"
	"log"

	"github.com/aomarai/PatchworkXIV/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Initialize OIDC middleware to verify tokens issued by Keycloak.
	// Services in Docker can reach Keycloak at http://keycloak:8080
	oidcMiddleware, err := middleware.OIDCMiddleware(context.Background(), "http://keycloak:8080/realms/xivmod", "xivmod-client")
	if err != nil {
		log.Fatalf("failed to initialize OIDC middleware: %v", err)
	}

	r.GET("/protected", oidcMiddleware, func(c *gin.Context) {
		claims, _ := c.Get("claims")
		c.JSON(200, gin.H{"message": "protected", "claims": claims})
	})

	r.Run(":8080")
}
