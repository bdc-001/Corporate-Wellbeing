package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// WebhookSignatureMiddleware validates webhook signatures
func WebhookSignatureMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secret == "" {
			c.Next()
			return
		}

		signature := c.GetHeader("X-Webhook-Signature")
		if signature == "" {
			c.JSON(401, gin.H{"error": "Missing webhook signature"})
			c.Abort()
			return
		}

		// Read request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to read request body"})
			c.Abort()
			return
		}

		// Restore body for handler
		c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

		// Verify signature
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expectedSignature := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			c.JSON(401, gin.H{"error": "Invalid webhook signature"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIKeyMiddleware validates API key authentication
func APIKeyMiddleware(apiKeySecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiKeySecret == "" {
			c.Next()
			return
		}

		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(401, gin.H{"error": "Missing API key"})
			c.Abort()
			return
		}

		// In production, validate against database
		// For now, simple comparison
		if apiKey != apiKeySecret {
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
