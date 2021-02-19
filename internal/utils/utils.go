package utils

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TimeLeftUntilMidnight() time.Duration {
	t := time.Now().AddDate(0, 0, 1)
	timezone, _ := time.LoadLocation("Asia/Seoul")
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, timezone)
	return time.Now().Sub(midnight)
}

func TodayTimeNoon() time.Time {
	t := time.Now().AddDate(0, 0, 1)
	timezone, _ := time.LoadLocation("Asia/Seoul")
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, timezone)
	return midnight
}

func CreateRandomString(length int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func CreateRandomNum(length int) string {
	chars := []rune("0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// HashAndSalt password
func HashAndSalt(origin string) string {
	pwd := []byte(origin)
	hash, hashAndSaltErr := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if hashAndSaltErr != nil {
		log.Println(hashAndSaltErr)
	}
	return string(hash)
}

// CheckPassword check password
func CheckPassword(normal, hashed string) bool {
	verifyHash := []byte(hashed)
	errFailedToCompareHAshAndPassword := bcrypt.CompareHashAndPassword(verifyHash, []byte(normal))
	if errFailedToCompareHAshAndPassword != nil {
		log.Println(errFailedToCompareHAshAndPassword)
		return false
	}
	return true
}

func ArrayHasItem(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func HasPermission(req []string, user []string) bool {
	if ok := ArrayHasItem(user, "admin"); ok {
		return true
	}
	for _, v := range req {
		if ok := ArrayHasItem(user, v); !ok {
			return false
		}
	}
	return true
}
