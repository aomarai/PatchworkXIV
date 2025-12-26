package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

// OIDCMiddleware returns a Gin middleware that verifies bearer tokens issued by the given issuer.
func OIDCMiddleware(ctx context.Context, issuer, clientID string) (gin.HandlerFunc, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := parts[1]
		idToken, err := verifier.Verify(ctx, token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var claims map[string]interface{}
		_ = idToken.Claims(&claims)
		c.Set("claims", claims)
		c.Next()
	}, nil
}
