package middleware

import (
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthLoginAdmin middleware checks login via Infomaniak OAuth2 and restricts to @sang.com users
func AuthLoginAdmin(c *fiber.Ctx) error {
	// Get token from header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	// Parse the token
	claims := &jwt.RegisteredClaims{}
	secretKey := []byte("your-secret-key") // Same key used for signing the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check if parsing token failed or token is invalid
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Check if the email belongs to the @sang.com domain
	email := claims.Subject
	if !strings.HasSuffix(email, "@sanggadee.com") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only @sanggadee.com emails are allowed for admin login"})
	}

	// Token is valid and email is authorized, proceed with request
	return c.Next()
}
