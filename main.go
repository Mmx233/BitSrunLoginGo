package main

import (
	"Mmx/Global"
	"Mmx/Request"
	"Mmx/Util"
	"encoding/json"
	"fmt"
)

func main() {
	fmt.Println("Step0: 检查状态…")
	G := Util.Config.Init()

	if Global.Config.Settings.QuitIfNetOk && Util.Checker.NetOk() {
		fmt.Println("网络正常，程序退出")
		return
	}

	fmt.Println("Step1: 正在获取客户端ip")
	{
		body, err := Request.Get(G.UrlLoginPage, nil)
		Util.ErrHandler(err)
		G.Ip, err = Util.GetIp(body)
		Util.ErrHandler(err)
	}
	fmt.Println("Step2: 正在获取Token")
	{
		data, err := Request.Get(G.UrlGetChallengeApi, map[string]string{
			"callback": "jsonp1583251661367",
			"username": G.Form.UserName,
			"ip":       G.Ip,
		})
		Util.ErrHandler(err)
		G.Token, err = Util.GetToken(data)
		Util.ErrHandler(err)
	}
	fmt.Println("Step3: 执行登录…")
	{
		info, err := json.Marshal(map[string]string{
			"username": G.Form.UserName,
			"password": G.Form.PassWord,
			"ip":       G.Ip,
			"acid":     G.Meta.Acid,
			"enc_ver":  G.Meta.Enc,
		})
		Util.ErrHandler(err)
		G.EncryptedInfo = "{SRBX1}" + Util.Base64(Util.XEncode(string(info), G.Token))
		G.Md5 = Util.Md5(G.Token)
		G.EncryptedMd5 = "{MD5}" + G.Md5

		var chkstr = G.Token + G.Form.UserName + G.Token + G.Md5
		chkstr += G.Token + G.Meta.Acid + G.Token + G.Ip
		chkstr += G.Token + G.Meta.N + G.Token + G.Meta.VType
		chkstr += G.Token + G.EncryptedInfo
		G.EncryptedChkstr = Util.Sha1(chkstr)

		res, err := Request.Get(G.UrlLoginApi, map[string]string{
			"callback":     "jQuery1124011576657442209481_1602812074032",
			"action":       "login",
			"username":     G.Form.UserName,
			"password":     G.EncryptedMd5,
			"ac_id":        G.Meta.Acid,
			"ip":           G.Ip,
			"info":         G.EncryptedInfo,
			"chksum":       G.EncryptedChkstr,
			"n":            G.Meta.N,
			"type":         G.Meta.VType,
			"os":           "Windows 10",
			"name":         "windows",
			"double_stack": "0",
			"_":            "1602812428675",
		})
		Util.ErrHandler(err)
		G.LoginResult, err = Util.GetResult(res)
		Util.ErrHandler(err)
		fmt.Println("登录结果: " + G.LoginResult)
		if Global.Config.Settings.DemoMode {
			fmt.Println(res)
		}
	}
}
