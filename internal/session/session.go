package session

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/osang-school/backend/graph/myerr"
	"github.com/osang-school/backend/internal/conf"
	"github.com/osang-school/backend/internal/db/redis"
	"github.com/osang-school/backend/internal/user"
	"github.com/osang-school/backend/internal/utils"
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
		Permission []string
		Role       user.Role
	}
)

const AccessTokenExp = time.Hour * 24 * 30 * 12

func CreateToken(id primitive.ObjectID, role user.Role, permission []string) (string, error) {
	randomStr := utils.CreateRandomString(8)
	claims := &Claims{Data: Data{
		ID:         id,
		SessionID:  randomStr,
		Permission: permission,
		Role:       role,
	},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExp).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(conf.Get().JWTSecret))

	redis.C.HSet("session:"+id.Hex(), randomStr, "Y")

	return tokenString, err
}

func ParseToken(tokenStr string) (*Data, error) {
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Get().JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, myerr.New(myerr.ErrAuth, "not valid token")
	}
	result, err := redis.C.HGet("session:"+claims.Data.ID.Hex(), claims.Data.SessionID).Result()
	if err != nil || result != "Y" {
		return nil, myerr.New(myerr.ErrAuth, "not valid token")
	}
	return &claims.Data, nil
}

func (d *Data) Expiry() error {
	return redis.C.HDel("session:"+d.ID.Hex(), d.SessionID).Err()
}
