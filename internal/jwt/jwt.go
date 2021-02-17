package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/osang-school/backend/internal/conf"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Claims struct {
		Data
		jwt.StandardClaims
	}

	Data struct {
		ID         primitive.ObjectID
		SessionID  string
		Permission *[]string
	}
)

const AccessTokenExp = time.Hour * 24 * 30 * 12

func createToken(id primitive.ObjectID, permission *[]string) (string, error) {
	claims := &Claims{
		Data: Data{
			ID:         id,
			Permission: permission,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExp).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(conf.Get().JWTSecret))

	return tokenString, err
}

func ParseToken(tokenStr string) (*Data, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Get().JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("Error Parse")
	}
	return &claims.Data, nil
}
