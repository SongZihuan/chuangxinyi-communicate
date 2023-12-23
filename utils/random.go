package utils

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/rand"
)

func GenerateRandomInt(min, max int) int64 {
	return int64(rand.GlobalRander.Intn(max-min+1) + min)
}

func GenerateUniqueNumber(length int) string {
	const numbers = "0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = numbers[rand.GlobalRander.Intn(len(numbers))]
	}

	return string(result)
}

func GenerateRandomText(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789?"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = charset[rand.GlobalRander.Intn(len(charset))]
	}

	return string(result)
}
