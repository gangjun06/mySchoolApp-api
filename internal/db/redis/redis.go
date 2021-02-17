package redis

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/osang-school/backend/internal/conf"
)

// C redis client
var C *redis.Client

func Init() {
	opt, err := redis.ParseURL(conf.Get().Redis.Addr)
	if err != nil {
		log.Fatal(err)
	}

	opt.Password = conf.Get().Redis.Pass

	C = redis.NewClient(opt)

	if result := C.Ping(); result.Err() != nil {
		panic(result.Err())
	}
}

func IsNil(err error) bool {
	if errors.Is(err, redis.Nil) {
		return true
	}
	return false
}

func SetObject(key string, value interface{}, expire time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return C.Set(key, p, expire).Err()
}

func GetObject(key string, dest interface{}) error {
	result := C.Get(key)
	if result.Err() != nil {
		return result.Err()
	}
	bytes, err := result.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dest)
}
