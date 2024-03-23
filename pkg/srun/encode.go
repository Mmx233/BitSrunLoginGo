package srun

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// Md5 编码
func Md5(token, password string) (string, error) {
	mac := hmac.New(md5.New, []byte(token))
	_, err := mac.Write([]byte(password))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// Sha1 编码
func Sha1(content string) string {
	h := sha1.New()
	h.Write([]byte(content))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x\n", bs)
}
