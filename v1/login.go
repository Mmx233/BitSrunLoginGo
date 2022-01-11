package BitSrun

import (
	"encoding/json"
	"errors"
	"github.com/Mmx233/BitSrunLoginGo/util"
	"github.com/Mmx233/BitSrunLoginGo/v1/transfer"
	"github.com/Mmx233/tool"
	"time"
)

func Login(c *srunTransfer.Login) error {
	util.Log.Demo = c.Demo
	util.Log.OutPut = c.OutPut

	G := util.GenerateLoginInfo(c.LoginInfo.Form, c.LoginInfo.Meta)
	if c.CheckNet {
		util.Log.Println("Step0: 检查状态…")
		if util.Checker.NetOk(c.Transport) {
			util.Log.Println("网络 ok")
			return nil
		}
	}

	util.Log.Println("Step1: 正在获取客户端ip")
	{
		if _, body, e := tool.HTTP.GetString(&tool.GetRequest{
			Url:       G.UrlLoginPage,
			Redirect:  true,
			Transport: c.Transport,
		}); e != nil {
			return e
		} else if G.Ip, e = util.GetIp(body); e != nil {
			return e
		}
	}

	util.Log.Println("Step2: 正在获取Token")
	{
		if _, data, e := tool.HTTP.GetString(&tool.GetRequest{
			Url: G.UrlGetChallengeApi,
			Query: map[string]interface{}{
				"callback": "jsonp1583251661367",
				"username": G.Form.UserName,
				"ip":       G.Ip,
			},
			Redirect:  true,
			Transport: c.Transport,
		}); e != nil {
			return e
		} else if G.Token, e = util.GetToken(data); e != nil {
			return e
		}
	}

	util.Log.Println("Step3: 执行登录…")
	{
		info, e := json.Marshal(map[string]string{
			"username": G.Form.UserName,
			"password": G.Form.PassWord,
			"ip":       G.Ip,
			"acid":     G.Meta.Acid,
			"enc_ver":  G.Meta.Enc,
		})
		if e != nil {
			return e
		}
		G.EncryptedInfo = "{SRBX1}" + util.Base64(util.XEncode(string(info), G.Token))
		G.Md5 = util.Md5(G.Token)
		G.EncryptedMd5 = "{MD5}" + G.Md5

		var chkstr = G.Token + G.Form.UserName + G.Token + G.Md5
		chkstr += G.Token + G.Meta.Acid + G.Token + G.Ip
		chkstr += G.Token + G.Meta.N + G.Token + G.Meta.Type
		chkstr += G.Token + G.EncryptedInfo
		G.EncryptedChkstr = util.Sha1(chkstr)

		if _, res, e := tool.HTTP.GetString(&tool.GetRequest{
			Url: G.UrlLoginApi,
			Query: map[string]interface{}{
				"callback":     "jQuery112401157665",
				"action":       "login",
				"username":     G.Form.UserName,
				"password":     G.EncryptedMd5,
				"ac_id":        G.Meta.Acid,
				"ip":           G.Ip,
				"info":         G.EncryptedInfo,
				"chksum":       G.EncryptedChkstr,
				"n":            G.Meta.N,
				"type":         G.Meta.Type,
				"os":           "Windows 10",
				"name":         "windows",
				"double_stack": 0,
				"_":            time.Now().UnixNano(),
			},
			Redirect:  true,
			Transport: c.Transport,
		}); e != nil {
			return e
		} else if G.LoginResult, e = util.GetResult(res); e != nil {
			return e
		} else {
			util.Log.Println("登录结果: " + G.LoginResult)
			if c.Demo {
				util.Log.Println(res)
			}
		}

		if G.LoginResult != "ok" {
			return errors.New(G.LoginResult)
		}
	}

	return nil
}
