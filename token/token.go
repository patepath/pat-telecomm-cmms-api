package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Username string `json:"username"`
	Position string `json:"position"`
	Role     uint8  `json:"role"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("aoise29387-12=123hr;1sfOGAS87vbfas49-*wpe98r1t23re123r102_(*&erq9ds;AOip293r")

func GenerateToken(username string, position string, role uint8) (string, error) {
	// Set custom and registered claims
	claims := &CustomClaims{
		Username: username,
		Position: position,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(14 * time.Hour)),
		},
	}

	// Create the token with the specified algorithm and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
	// Parse the token and validate its signature and claims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Cast the token claims to the custom claims type
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
