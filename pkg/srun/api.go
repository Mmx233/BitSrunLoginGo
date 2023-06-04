package srun

import (
	"encoding/json"
	"fmt"
	"github.com/Mmx233/tool"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Api struct {
	BaseUrl string
	Client  *http.Client
	// 禁用自动重定向
	NoDirect *http.Client
}

func (a *Api) Init(https bool, domain string, client *http.Client) {
	a.BaseUrl = "http"
	if https {
		a.BaseUrl += "s"
	}
	a.BaseUrl = a.BaseUrl + "://" + domain + "/"

	// 初始化 http client
	a.Client = client
	a.NoDirect = &(*client)
	a.NoDirect.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
}

func (a *Api) request(path string, query map[string]interface{}) (map[string]interface{}, error) {
	log.Debugln("HTTP GET ", a.BaseUrl+path)
	timestamp := fmt.Sprint(time.Now().UnixNano())
	callback := "jQuery" + timestamp
	if query == nil {
		query = make(map[string]interface{}, 2)
	}
	query["callback"] = callback
	query["_"] = timestamp
	httpTool := tool.NewHttpTool(a.Client)
	req, e := httpTool.GenReq("GET", &tool.DoHttpReq{
		Url:   a.BaseUrl + path,
		Query: query,
	})
	if e != nil {
		log.Debugln(e)
		return nil, e
	}

	resp, e := httpTool.Client.Do(req)
	if e != nil {
		log.Debugln(e)
		return nil, e
	}
	defer resp.Body.Close()

	data, e := io.ReadAll(resp.Body)
	if e != nil {
		log.Debugln(e)
		return nil, e
	}
	res := string(data)

	log.Debugln(res)
	res = strings.TrimPrefix(res, callback+"(")
	res = strings.TrimSuffix(res, ")")

	var r map[string]interface{}
	return r, json.Unmarshal([]byte(res), &r)
}

func (a *Api) GetUserInfo() (map[string]interface{}, error) {
	return a.request("cgi-bin/rad_user_info", nil)
}

func (a *Api) DetectAcid() (string, error) {
	addr := a.BaseUrl
	for {
		log.Debugln("HTTP GET ", addr)
		res, e := a.NoDirect.Get(addr)
		if e != nil {
			return "", e
		}
		_ = res.Body.Close()
		loc := res.Header.Get("location")
		if res.StatusCode == 302 && loc != "" {
			if strings.HasPrefix(loc, "/") {
				addr = a.BaseUrl + strings.TrimPrefix(loc, "/")
			} else {
				addr = loc
			}

			var u *url.URL
			u, e = url.Parse(addr)
			if e != nil {
				return "", e
			}
			acid := u.Query().Get(`ac_id`)
			if acid != "" {
				return acid, nil
			}

			continue
		}
		break
	}
	return "", ErrAcidCannotFound
}

func (a *Api) Login(
	Username,
	Password,
	AcID,
	Ip,
	Info,
	ChkSum,
	N,
	Type string,
) (map[string]interface{}, error) {
	return a.request(
		"cgi-bin/srun_portal",
		map[string]interface{}{
			"action":       "login",
			"username":     Username,
			"password":     Password,
			"ac_id":        AcID,
			"ip":           Ip,
			"info":         Info,
			"chksum":       ChkSum,
			"n":            N,
			"type":         Type,
			"os":           "Windows 10",
			"name":         "windows",
			"double_stack": 0,
		})
}

func (a *Api) GetChallenge(username, ip string) (map[string]interface{}, error) {
	return a.request(
		"cgi-bin/get_challenge",
		map[string]interface{}{
			"username": username,
			"ip":       ip,
		})
}
