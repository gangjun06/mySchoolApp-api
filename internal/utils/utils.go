package utils

import (
	"math/rand"
	"strings"
	"time"
)

func TimeLeftUntilMidnight() time.Duration {
	t := time.Now().AddDate(0, 0, 1)
	timezone, _ := time.LoadLocation("Asia/Seoul")
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, timezone)
	return time.Now().Sub(midnight)
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
