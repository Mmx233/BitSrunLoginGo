package BitSrun

import (
	"encoding/json"
	"fmt"
	"github.com/Mmx233/BitSrunLoginGo/util"
	"github.com/Mmx233/tool"
	"net/http"
	"strings"
	"time"
)

type SrunApi struct {
	BaseUrl string
	Client  *http.Client
}

func (a *SrunApi) request(path string, query map[string]interface{}) (map[string]interface{}, error) {
	util.Log.Debug("HTTP GET ", a.BaseUrl+path)
	timestamp := fmt.Sprint(time.Now().UnixNano())
	callback := "jQuery" + timestamp
	if query == nil {
		query = make(map[string]interface{}, 2)
	}
	query["callback"] = callback
	query["_"] = timestamp
	_, res, e := tool.NewHttpTool(a.Client).GetString(&tool.DoHttpReq{
		Url:   a.BaseUrl + path,
		Query: query,
	})
	if e != nil {
		util.Log.Debug(e)
		return nil, e
	}

	util.Log.Debug(res)
	res = strings.TrimPrefix(res, callback+"(")
	res = strings.TrimSuffix(res, ")")

	var r map[string]interface{}
	return r, json.Unmarshal([]byte(res), &r)
}

func (a *SrunApi) GetUserInfo() (map[string]interface{}, error) {
	return a.request("cgi-bin/rad_user_info", nil)
}

func (a *SrunApi) Login(
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

func (a *SrunApi) GetChallenge(username, ip string) (map[string]interface{}, error) {
	return a.request(
		"cgi-bin/get_challenge",
		map[string]interface{}{
			"username": username,
			"ip":       ip,
		})
}
