package utils

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/spaolacci/murmur3"
)

func HashSHA256WithBase62(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hash := hasher.Sum(nil)
	return EncodeToBase62(hash)
}

func HashSHA256(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

// 编码用的字符集
func MurmurHashWithBase62(s string) string {
	// 使用MurmurHash对字符串进行哈希
	h := murmur3.New64()
	_, _ = h.Write([]byte(s))
	hash := h.Sum64()

	// 将哈希值转换为62进制编码
	return EncodeBase62(hash)
}

func HashSHA1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
