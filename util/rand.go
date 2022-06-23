package util

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

// RandInt return [0-n)
func RandInt(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

// RandIntWithoutZero return [1-n]
func RandIntWithoutZero(n int) int {
	return RandInt(n) + 1
}

// GenerateRandomAlpha rand
func GenerateRandomAlpha(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterRunes))]
	}
	return string(b)
}

// GenerateUuid .
func GenerateUuid() string {
	return uuid.New().String()
}
