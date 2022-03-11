package util

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
)

// Md5 编码
func Md5(content string) string {
	w := md5.New()
	_, _ = io.WriteString(w, content)
	return fmt.Sprintf("%x", w.Sum(nil))
}

// Sha1 编码
func Sha1(content string) string {
	h := sha1.New()
	h.Write([]byte(content))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x\n", bs)
}
