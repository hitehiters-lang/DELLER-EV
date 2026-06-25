package pkgauth

import (
	"deller/internal/config"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JWTSecret = config.GetJwt()
var jwtKey = []byte(JWTSecret)

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uuid.UUID) (string, error) {
	claims := Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2400 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	cleanToken := strings.TrimPrefix(tokenString, "Bearer ")
	cleanToken = strings.TrimSpace(cleanToken)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cleanToken, claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtKey, nil
		})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
