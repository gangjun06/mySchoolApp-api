package utils

import (
	"fmt"
	"time"

	"github.com/gangjun06/mySchoolApp-api/models"

	"github.com/dgrijalva/jwt-go"
)

func GetJwtToken(id uint) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &models.Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, errStringToSignedString := token.SignedString([]byte(GetConfig().Server.JWTSecret))

	if errStringToSignedString != nil {
		fmt.Println(errStringToSignedString)
		return "", fmt.Errorf("token signed Error")
	}
	return tokenString, nil
}
