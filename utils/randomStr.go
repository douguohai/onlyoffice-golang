package utils

import (
	"math/rand"
	"time"
)

// RandomStringGenerator 随机字符串生成器
type RandomStringGenerator struct {
	// 字符集
	charset string

	// 字符串长度
	length int
}

// NewRandomStringGenerator 创建随机字符串生成器
func NewRandomStringGenerator(charset string, length int) *RandomStringGenerator {
	return &RandomStringGenerator{
		charset: charset,
		length:  length,
	}
}

// Generate 生成随机字符串
func (g *RandomStringGenerator) Generate() string {
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, g.length)
	for i := 0; i < g.length; i++ {
		result[i] = g.charset[rand.Intn(len(g.charset))]
	}

	return string(result)
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	generator := NewRandomStringGenerator("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", length)
	return generator.Generate()
}
