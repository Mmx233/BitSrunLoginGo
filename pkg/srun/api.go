package srun

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mmx233/tool"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
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

	CustomHeader map[string]interface{}
}

type ApiConfig struct {
	Https        bool
	Domain       string
	Client       *http.Client
	CustomHeader map[string]interface{}
}

func (a *Api) Init(conf *ApiConfig) {
	a.BaseUrl = "http"
	if conf.Https {
		a.BaseUrl += "s"
	}
	a.BaseUrl = a.BaseUrl + "://" + conf.Domain + "/"

	a.CustomHeader = conf.CustomHeader

	// 初始化 http client
	a.Client = conf.Client
	copyClient := *conf.Client
	a.NoDirect = &copyClient
	a.NoDirect.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
}

func (a *Api) request(path string, query map[string]interface{}) (map[string]interface{}, error) {
	log.Debugln("HTTP GET ", a.BaseUrl+path)
	callback := fmt.Sprintf("jQuery%s_%d", tool.NewRand(rand.NewSource(time.Now().UnixNano())).WithLetters("123456789").String(21), time.Now().UnixMilli())
	if query == nil {
		query = make(map[string]interface{}, 2)
	}
	query["callback"] = callback
	query["_"] = fmt.Sprint(time.Now().UnixMilli())
	httpTool := tool.NewHttpTool(a.Client)
	req, err := httpTool.GenReq("GET", &tool.DoHttpReq{
		Url:    a.BaseUrl + path,
		Query:  query,
		Header: a.CustomHeader,
	})
	if err != nil {
		log.Debugln(err)
		return nil, err
	}

	resp, err := httpTool.Client.Do(req)
	if err != nil {
		log.Debugln(err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Debugln(err)
		return nil, err
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

func (a *Api) FollowRedirect(addr *url.URL, onNext func(addr *url.URL) error) (*url.URL, error) {
	addrCopy := *addr
	addr = &addrCopy
	for {
		log.Debugln("HTTP GET ", addr)
		req, err := http.NewRequest("GET", addr.String(), nil)
		if err != nil {
			return nil, err
		}
		for k, v := range a.CustomHeader {
			req.Header.Set(k, fmt.Sprint(v))
		}
		res, err := a.NoDirect.Do(req)
		if err != nil {
			return nil, err
		}
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
		loc := res.Header.Get("location")
		if res.StatusCode < 300 {
			break
		} else if res.StatusCode < 400 {
			if loc == "" {
				return nil, errors.New("目标跳转地址缺失")
			}
			if strings.HasPrefix(loc, "/") {
				addr.Path = strings.TrimPrefix(loc, "/")
			} else {
				addr, err = url.Parse(loc)
				if err != nil {
					return nil, err
				}
			}
			if err = onNext(addr); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("server return http status %d", res.StatusCode)
		}
	}
	return addr, nil
}

func (a *Api) searchAcid(query url.Values) (string, bool) {
	addr := query.Get(`ac_id`)
	return addr, addr != ""
}

// DetectAcid err 为 nil 时 acid 一定存在
func (a *Api) DetectAcid() (string, error) {
	baseUrl, err := url.Parse(a.BaseUrl)
	if err != nil {
		return "", err
	}

	var AcidFound = errors.New("acid found")
	var acid string
	_, err = a.FollowRedirect(baseUrl, func(addr *url.URL) error {
		var ok bool
		acid, ok = a.searchAcid(addr.Query())
		if ok {
			return AcidFound
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, AcidFound) {
			return acid, nil
		}
		return "", err
	}
	return "", ErrAcidCannotFound
}

// Reality acid 可能为空字符串
func (a *Api) Reality(addr string, getAcid bool) (acid string, online bool, err error) {
	startUrl, err := url.Parse(addr)
	if err != nil {
		return "", false, err
	}
	var AlreadyOnline = errors.New("already online")
	var finalUrl *url.URL
	finalUrl, err = a.FollowRedirect(startUrl, func(addr *url.URL) error {
		// 任一跳转没有跳出初始域名说明已经在线
		if addr.Host == startUrl.Host {
			return AlreadyOnline
		}
		if getAcid {
			acid, _ = a.searchAcid(addr.Query())
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, AlreadyOnline) {
			online = true
			err = nil
			return
		}
		return
	}
	online = finalUrl.Host == startUrl.Host
	return
}

type LoginRequest struct {
	Username    string
	Password    string
	AcID        string
	IP          string
	Info        string
	ChkSum      string
	N           string
	Type        string
	OS          string
	Name        string
	DoubleStack string
}

func (a *Api) Login(req *LoginRequest) (map[string]interface{}, error) {
	return a.request(
		"cgi-bin/srun_portal",
		map[string]interface{}{
			"action":       "login",
			"username":     req.Username,
			"password":     req.Password,
			"ac_id":        req.AcID,
			"ip":           req.IP,
			"info":         req.Info,
			"chksum":       req.ChkSum,
			"n":            req.N,
			"type":         req.Type,
			"os":           req.OS,
			"name":         req.Name,
			"double_stack": req.DoubleStack,
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
