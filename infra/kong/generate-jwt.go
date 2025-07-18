package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Kong JWT Configuration
const (
	JWTKey    = "bitzap-key"
	JWTSecret = "bitzap-secret-2025"
)

// JWTClaims represents the JWT payload
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// generateJWTToken creates a JWT token for Kong authentication
func generateJWTToken(userID string, expiresInHours int) (string, error) {
	now := time.Now()

	claims := JWTClaims{
		UserID:   userID,
		Username: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    JWTKey, // Kong key
			Subject:   userID, // User ID
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiresInHours) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// decodeJWTToken decodes and verifies JWT token
func decodeJWTToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func main() {
	fmt.Println("🔑 Kong JWT Token Generator (Go Version)")
	fmt.Println("=" + fmt.Sprintf("%40s", ""))

	// Generate token
	userID := "testuser"
	expiresInHours := 24

	token, err := generateJWTToken(userID, expiresInHours)
	if err != nil {
		log.Fatalf("Error generating token: %v", err)
	}

	fmt.Printf("Generated JWT Token:\n%s\n\n", token)

	// Decode token to show payload
	claims, err := decodeJWTToken(token)
	if err != nil {
		log.Fatalf("Error decoding token: %v", err)
	}

	// Pretty print claims
	claimsJSON, _ := json.MarshalIndent(claims, "", "  ")
	fmt.Printf("Token Payload:\n%s\n\n", string(claimsJSON))

	fmt.Printf("Usage Example:\n")
	fmt.Printf("curl -H 'Authorization: Bearer %s' http://localhost:8000/api/v1/shorten\n\n", token)

	fmt.Printf("Token expires at: %s\n", claims.RegisteredClaims.ExpiresAt.Time.Format("2006-01-02 15:04:05"))
}
