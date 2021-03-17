package main

import (
	"Mmx/Global"
	"Mmx/Request"
	"Mmx/Util"
	"encoding/json"
	"fmt"
)

func main() {
	G := Util.Config.Init()

	if Global.Config.Settings.QuitIfNetOk && Util.Checker.NetOk() {
		fmt.Println("There's no need to login")
		return
	}

	fmt.Println("Step1: Get local ip returned from srun server.")
	{
		body, err := Request.Get(G.UrlLoginPage, nil)
		Util.ErrHandler(err)
		G.Ip, err = Util.GetIp(body)
		Util.ErrHandler(err)
	}
	fmt.Println("Step2: Get token by resolving challenge result.")
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
	fmt.Println("Step3: Login and resolve response.")
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

		var chkstr string
		chkstr = G.Token + G.Form.UserName
		chkstr += G.Token + G.Md5
		chkstr += G.Token + G.Meta.Acid
		chkstr += G.Token + G.Ip
		chkstr += G.Token + G.Meta.N
		chkstr += G.Token + G.Meta.VType
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
	}
	fmt.Println("The login result is: " + G.LoginResult)
}
