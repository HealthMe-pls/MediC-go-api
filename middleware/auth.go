package middleware

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/HealthMe-pls/medic-go-api/database"
)

// Secret key for JWT
var jwtSecret = []byte("your_secret_key")

// AuthLogin middleware checks for a valid JWT token and ensures it is not blacklisted
func AuthLogin(c *fiber.Ctx) error {
	// Get token from header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	// Check if token is blacklisted in Redis
	exists, err := database.RedisClient.Get(context.Background(), tokenString).Result()
	if err == nil && exists == "blacklisted" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token has been revoked"})
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Token is valid, allow the request to continue
	return c.Next()
}
