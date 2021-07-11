package util

import (
	"Mmx/global"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"time"
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
	return Search("id=\"user_ip\" value=\"(.*?)\"", body)
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

func ErrHandler(err error) {
	if !global.Status.Output {
		return
	}
	Log.Println("运行出错，状态异常")
	if global.Config.Settings.DemoMode {
		Log.Fatalln(err)
	}
	os.Exit(1)
}

func NetDailEr() func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: 5 * time.Second,
		}
		return d.DialContext(ctx, "udp", global.Config.Settings.Dns+":53")
	}
}
