package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"github.com/gorilla/mux"
)

// Define the expected structure for login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Define the response structure
type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Handle Login request
func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into LoginRequest
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Now that you have the email and password, authenticate them
	// You can implement actual authentication logic here (e.g., verify with Infomaniak)
	token, err := authenticateUser(loginReq.Email, loginReq.Password)
	if err != nil {
		// If authentication fails, send an error response
		resp := LoginResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// If authentication is successful, return a success response with the token
	resp := LoginResponse{
		Success: true,
		Token:   token,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Simulate user authentication (replace with actual logic)
func authenticateUser(email, password string) (string, error) {
	// Placeholder logic: Check the credentials and return a mock token
	// Replace with Infomaniak OAuth2 or your own authentication mechanism
	if email == "admin@sang.com" && password == "password123" {
		return "mock-jwt-token", nil // Return a mock token on success
	}
	return "", fmt.Errorf("invalid credentials")
}

func main() {
	r := mux.NewRouter()

	// Define the login route
	r.HandleFunc("/api/auth/admin", LoginAdmin).Methods("POST")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Example function to generate a JWT token
func GenerateJWT(email string) (string, error) {
	// Define token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)
	
	// Create the JWT claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Issuer:    "your-app-name", // Replace with your app name or domain
		Subject:   email,
	}

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Replace with your secret key
	secretKey := []byte("your-secret-key")
	
	// Sign the token and return it
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
