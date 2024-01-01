package srun

import (
	"encoding/json"
	"fmt"
	"github.com/Mmx233/tool"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
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
	log.Debugln("HTTP GET", a.BaseUrl+path)
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
