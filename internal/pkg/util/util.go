package util

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

// GenerateUUID 生成 UUID（不含连字符）
func GenerateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return strings.ReplaceAll(hex.EncodeToString(b), "-", "")
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// TruncateString 截断字符串
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// SanitizeString 清理字符串（去除前后空格）
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}
