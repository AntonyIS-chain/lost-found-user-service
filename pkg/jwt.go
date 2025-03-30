package pkg

import (
	"errors"
	"time"

	"github.com/AntonyIS-chain/lost-found-user-service/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	accessTokenTTL  = time.Minute * 15   // Access token expires in 15 minutes
	refreshTokenTTL = time.Hour * 24 * 7 // Refresh token expires in 7 days
)

// GenerateToken creates a new JWT access token
func GenerateToken(id string) (string, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"id": id,
		"exp":   time.Now().Add(accessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// ✅ Convert SECRET_KEY to []byte
	return token.SignedString([]byte(cfg.SECRET_KEY))
}

// GenerateRefreshToken creates a new JWT refresh token
func GenerateRefreshToken(id string) (string, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"id": id,
		"exp":   time.Now().Add(refreshTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// ✅ Convert SECRET_KEY to []byte
	return token.SignedString([]byte(cfg.SECRET_KEY))
}

// ValidateRefreshToken validates a refresh token and returns the id
func ValidateRefreshToken(tokenString string) (string, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.SECRET_KEY), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired refresh token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid refresh token payload")
	}

	// Extract id
	id, ok := claims["id"].(string)
	if !ok {
		return "", errors.New("invalid id in refresh token")
	}

	// Extract expiration and check if it's expired
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", errors.New("invalid expiration in refresh token")
	}

	// Convert `exp` to `time.Time` and validate
	if time.Now().Unix() > int64(exp) {
		return "", errors.New("refresh token has expired")
	}

	return id, nil
}
