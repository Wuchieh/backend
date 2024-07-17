package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Sha256 將接收到的字串合併並使用sha256加密
func Sha256(s ...string) string {
	str := strings.Join(s, "")

	// 使用 SHA-256 加密
	harsher := sha256.New()
	harsher.Write([]byte(str))
	hashBytes := harsher.Sum(nil)

	// 將加密後的結果轉換成十六進位字串
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
