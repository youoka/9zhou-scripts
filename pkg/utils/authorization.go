package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

// CheckIsValid
func CheckIsValid(encKey, shop, reclaim string) bool {
	authorization, err := DecryptAuthorization(encKey, "9zhou-scripts-default-key-32-chars")
	if err != nil {
		return false
	}
	split := strings.Split(authorization, "|")
	if len(split) < 2 {
		return false
	}
	if split[0] != shop || split[1] != reclaim {
		return false
	}
	date := extractDate(authorization)
	// 检查日期是否小于当前日期，先解析string成time
	currentTime := time.Now()
	expirationTime, _ := time.Parse("2006-01-02", date)
	return currentTime.Before(expirationTime) || currentTime.Equal(expirationTime)
}

// 授权码加密
// 明文格式: 商城号核销号到期时间(2025-10-12)
func HashAuthorization(authCode string) (string, error) {
	// 验证授权码格式
	if !IsValidAuthCodeFormat(authCode) {
		return "", &AuthorizationError{"invalid authorization code format"}
	}

	// 使用默认密钥进行加密
	return EncryptAuthorization(authCode, "9zhou-scripts-default-key-32-chars")
}

// 授权码验证
func CheckAuthorization(authCode, encrypted string) bool {
	// 验证授权码格式
	if !IsValidAuthCodeFormat(authCode) {
		return false
	}

	// 解密已加密的授权码
	decrypted, err := DecryptAuthorization(encrypted, "9zhou-scripts-default-key-32-chars")
	if err != nil {
		return false
	}

	// 提取并比较日期
	decryptedDate := extractDate(decrypted)
	authCodeDate := extractDate(authCode)

	// 检查日期是否一致且未过期
	if decryptedDate != authCodeDate {
		return false
	}

	currentTime := time.Now()
	expirationTime, _ := time.Parse("2006-01-02", authCodeDate)

	return currentTime.Before(expirationTime) || currentTime.Equal(expirationTime)
}

// EncryptAuthorization 使用AES加密授权码
func EncryptAuthorization(plaintext, key string) (string, error) {
	// 创建AES加密块
	block, err := aes.NewCipher([]byte(key)[:32])
	if err != nil {
		return "", err
	}

	// 将明文转换为字节数组
	plaintextBytes := []byte(plaintext)

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 创建随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, plaintextBytes, nil)

	// 将密文编码为base64字符串
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAuthorization 使用AES解密授权码
func DecryptAuthorization(encrypted, key string) (string, error) {
	// 解码base64字符串
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	// 创建AES加密块
	block, err := aes.NewCipher([]byte(key)[:32])
	if err != nil {
		return "", err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 获取nonce大小
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// 分离nonce和密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// 验证授权码格式
// 格式: 商城号核销号到期时间(2025-10-12)
func IsValidAuthCodeFormat(authCode string) bool {
	// 检查是否以日期结尾，格式为 (YYYY-MM-DD)
	matched, _ := regexp.MatchString(`^.+\|\(.+\)\(\d{4}-\d{2}-\d{2}\)$`, authCode)
	if !matched {
		// 兼容旧格式：没有商城号和核销号前缀的情况
		matched, _ = regexp.MatchString(`^.+\(\d{4}-\d{2}-\d{2}\)$`, authCode)
		if !matched {
			return false
		}
	}

	// 提取日期部分进行验证
	datePart := extractDate(authCode)
	if datePart == "" {
		return false
	}

	_, err := time.Parse("2006-01-02", datePart)
	return err == nil
}

// extractDate 从授权码中提取日期部分
func extractDate(authCode string) string {
	start := strings.LastIndex(authCode, "(")
	end := strings.LastIndex(authCode, ")")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return authCode[start+1 : end]
}

// ExtractDate 公开版本的extractDate函数
func ExtractDate(authCode string) string {
	return extractDate(authCode)
}

// AuthorizationError represents an authorization-related error
type AuthorizationError struct {
	msg string
}

func (e *AuthorizationError) Error() string {
	return e.msg
}
