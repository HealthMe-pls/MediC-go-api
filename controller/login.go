package controller

import (
	"errors"
	"io"
	"time"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)


// Secret key for JWT
var jwtSecret = []byte("your_secret_key")
var secretKey = []byte("mysecretencryptionkey123") // Must be 16, 24, or 32 bytes


// HashPassword hashes a password
// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	return string(bytes), err
// }
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
// // CheckPassword compares a hashed password with plain text
// func CheckPassword(hashedPassword, password string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
// 	return err == nil
// }
func CheckPassword(encryptedPassword, inputPassword string) bool {
	// Decrypt the stored encrypted password
	decryptedPassword, err := DecryptPassword(encryptedPassword)
	if err != nil {
		return false // Return false if decryption fails
	}

	// Compare decrypted password with the user input
	return decryptedPassword == inputPassword
}


// EncryptPassword encrypts a password using AES
func EncryptPassword(password string) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(password))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(password))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPassword decrypts an encrypted password
func DecryptPassword(encrypted string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
// Register a new Entrepreneur
func Register(db *gorm.DB, c *fiber.Ctx) error {
	var input model.Entrepreneur

	// Parse request body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Hash and encrypt the password
	// hashedPassword, err := HashPassword(input.Password)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	// }
	encryptedPassword, err := EncryptPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encrypt password"})
	}

	// Create a new Entrepreneur
	Entrepreneur := model.Entrepreneur{
		Username: input.Username,
		Password: encryptedPassword, // Store hashed password
	}

	// Save to database
	if err := db.Create(&Entrepreneur).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":            "User registered successfully",
		"username":           Entrepreneur.Username,
		"password_encrypted": encryptedPassword, // Show encrypted password for debugging (Remove in production)
	})
}

// GenerateJWT generates a JWT token
func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Login function using Entrepreneur model directly
func Login(db *gorm.DB, c *fiber.Ctx) error {
	var input model.Entrepreneur
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Find entrepreneur by username
	var entrepreneur model.Entrepreneur
	if err := db.Where("username = ?", input.Username).First(&entrepreneur).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Check password
	if !CheckPassword(entrepreneur.Password, input.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT
	token, err := GenerateJWT(entrepreneur.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

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
// Logout function to blacklist JWT token
func Logout(c *fiber.Ctx) error {
	// Get token from header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Store the token in Redis with an expiration time
	err := database.RedisClient.Set(context.Background(), tokenString, "blacklisted", 24*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to blacklist token"})
	}

	return c.JSON(fiber.Map{"message": "Successfully logged out"})
}
// GetEntrepreneurWithPassword fetches entrepreneur details including the hashed password
func GetEntrepreneurWithPassword(db *gorm.DB, c *fiber.Ctx) error {
	username := c.Params("username") // Get username from URL params

	var entrepreneur model.Entrepreneur
	if err := db.Where("username = ?", username).First(&entrepreneur).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Entrepreneur not found"})
	}

	// WARNING: Returning hashed passwords is a security risk!
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"username": entrepreneur.Username,
		"password": entrepreneur.Password, // Hashed password
		"message":  "Entrepreneur found",
	})
}


func GetAllEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
	var entrepreneur []model.Entrepreneur

	// Fetch all entrepreneur
	if err := db.Find(&entrepreneur).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve entrepreneur",
			"details": err.Error(),
		})
	}

	// Decrypt passwords
	for i := range entrepreneur {
		decryptedPassword, err := DecryptPassword(entrepreneur[i].Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to decrypt password",
				"details": err.Error(),
			})
		}
		entrepreneur[i].Password = decryptedPassword
	}

	return c.Status(fiber.StatusOK).JSON(entrepreneur)
}
