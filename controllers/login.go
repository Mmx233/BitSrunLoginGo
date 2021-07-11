package controllers

import (
	"Mmx/global"
	"Mmx/request"
	"Mmx/util"
	"encoding/json"
	"fmt"
	"time"
)

func Login(output bool) {
	global.Status.Output = output
	util.Log.Println("Step0: 检查状态…")
	G := global.Config.Generate()

	if !global.Status.Daemon && global.Config.Settings.QuitIfNetOk && util.Checker.NetOk() {
		util.Log.Println("网络正常，程序退出")
		return
	}

	util.Log.Println("Step1: 正在获取客户端ip")
	{
		body, err := request.Get(G.UrlLoginPage, nil)
		if err != nil {
			util.ErrHandler(err)
			return
		}
		G.Ip, err = util.GetIp(body)
		if err != nil {
			util.ErrHandler(err)
			return
		}
	}
	util.Log.Println("Step2: 正在获取Token")
	{
		data, err := request.Get(G.UrlGetChallengeApi, map[string]string{
			"callback": "jsonp1583251661367",
			"username": G.Form.UserName,
			"ip":       G.Ip,
		})
		if err != nil {
			util.ErrHandler(err)
			return
		}
		G.Token, err = util.GetToken(data)
		if err != nil {
			util.ErrHandler(err)
			return
		}
	}
	util.Log.Println("Step3: 执行登录…")
	{
		info, err := json.Marshal(map[string]string{
			"username": G.Form.UserName,
			"password": G.Form.PassWord,
			"ip":       G.Ip,
			"acid":     G.Meta.Acid,
			"enc_ver":  G.Meta.Enc,
		})
		if err != nil {
			util.ErrHandler(err)
			return
		}
		G.EncryptedInfo = "{SRBX1}" + util.Base64(util.XEncode(string(info), G.Token))
		G.Md5 = util.Md5(G.Token)
		G.EncryptedMd5 = "{MD5}" + G.Md5

		var chkstr = G.Token + G.Form.UserName + G.Token + G.Md5
		chkstr += G.Token + G.Meta.Acid + G.Token + G.Ip
		chkstr += G.Token + G.Meta.N + G.Token + G.Meta.Type
		chkstr += G.Token + G.EncryptedInfo
		G.EncryptedChkstr = util.Sha1(chkstr)

		res, err := request.Get(G.UrlLoginApi, map[string]string{
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
			"double_stack": "0",
			"_":            fmt.Sprint(time.Now().UnixNano()),
		})
		if err != nil {
			util.ErrHandler(err)
			return
		}
		G.LoginResult, err = util.GetResult(res)
		if err != nil {
			util.ErrHandler(err)
			return
		}
		util.Log.Println("登录结果: " + G.LoginResult)
		if global.Config.Settings.DemoMode {
			util.Log.Println(res)
		}
	}
}
