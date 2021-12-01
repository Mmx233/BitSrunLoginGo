package util

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"regexp"
)

func Search(reg string, content string) (string, error) {
	r := regexp.MustCompile(reg)
	if r == nil {
		return "", errors.New("解析正则表达式失败")
	}
	if s := r.FindStringSubmatch(content); len(s) < 2 {
		return "", errors.New("无匹配")
	} else {
		return s[1], nil
	}
}

func GetIp(body string) (string, error) {
	//判断原正则是否有匹配，如果无就使用新正则尝试
	if ip, e := Search("id=\"user_ip\" value=\"(.*?)\"", body); e == nil {
		return ip, nil
	}
	return Search("ip     : \"(.*?)\"", body)
}

func GetToken(body string) (string, error) {
	return Search("\"challenge\":\"(.*?)\"", body)
}

func GetResult(body string) (string, error) {
	return Search("\"error\":\"(.+?)\"", body)
}

func Md5(content string) string {
	w := md5.New()
	_, _ = io.WriteString(w, content)
	return fmt.Sprintf("%x", w.Sum(nil))
}

func Sha1(content string) string {
	h := sha1.New()
	h.Write([]byte(content))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x\n", bs)
}
