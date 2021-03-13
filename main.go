package main

import (
	"Mmx/Request"
	"Mmx/Util"
	"encoding/json"
	"fmt"
	"os"
)

type LoginForm struct {
	UserName string
	PassWord string
}

type LoginInfo struct {
	UrlLoginPage       string
	UrlGetChallengeApi string
	UrlLoginApi        string
	N                  string
	VType              string
	Acid               string
	Enc                string

	Ip              string
	Token           string
	EncryptedInfo   string
	Md5             string
	EncryptedMd5    string
	EncryptedChkstr string
	LoginResult     string
	Form            *LoginForm
}

func Generate(Domain string, Username string, Password string) *LoginInfo {
	return &LoginInfo{
		UrlLoginPage:       "http://" + Domain + "/srun_portal_success?ac_id=5&theme=basic1",
		UrlGetChallengeApi: "http://" + Domain + "/cgi-bin/get_challenge",
		UrlLoginApi:        "http://" + Domain + "/cgi-bin/srun_portal",
		N:                  "200",
		VType:              "1",
		Acid:               "5",
		Enc:                "srun_bx1",
		Form: &LoginForm{
			UserName: Username + "@cmcc",
			PassWord: Password,
		},
	}
}

func ErrHandler(err error) {
	if err != nil {
		defer os.Exit(3)
		fmt.Println("Error")
		panic(err)
	}
}

func main() {
	G := Generate(
		"", //登录地址域名或ip
		"", //账号
		"", //密码
	)
	fmt.Println("Step1: Get local ip returned from srun server.")
	{
		body, err := Request.Get(G.UrlLoginPage, nil)
		ErrHandler(err)
		G.Ip, err = Util.GetIp(body)
		ErrHandler(err)
	}
	fmt.Println("Step2: Get token by resolving challenge result.")
	{
		data, err := Request.Get(G.UrlGetChallengeApi, map[string]string{
			"callback": "jsonp1583251661367",
			"username": G.Form.UserName,
			"ip":       G.Ip,
		})
		ErrHandler(err)
		G.Token, err = Util.GetToken(data)
		ErrHandler(err)
	}
	fmt.Println("Step3: Loggin and resolve response.")
	{
		info, err := json.Marshal(map[string]string{
			"username": G.Form.UserName,
			"password": G.Form.PassWord,
			"ip":       G.Ip,
			"acid":     G.Acid,
			"enc_ver":  G.Enc,
		})
		ErrHandler(err)
		G.EncryptedInfo = "{SRBX1}" + Util.Base64(Util.XEncode(string(info), G.Token))
		G.Md5=Util.Md5(G.Token)
		G.EncryptedMd5 = "{MD5}"+G.Md5

		var chkstr string
		chkstr = G.Token + G.Form.UserName
		chkstr += G.Token + G.Md5
		chkstr += G.Token + G.Acid
		chkstr += G.Token + G.Ip
		chkstr += G.Token + G.N
		chkstr += G.Token + G.VType
		chkstr += G.Token+ G.EncryptedInfo
		G.EncryptedChkstr=Util.Sha1(chkstr)

		res,err:=Request.Get(G.UrlLoginApi, map[string]string{
			"callback": "jQuery1124011576657442209481_1602812074032",
			"action":"login",
			"username": G.Form.UserName,
			"password": G.EncryptedMd5,
			"ac_id": G.Acid,
			"ip": G.Ip,
			"info": G.EncryptedInfo,
			"chksum": G.EncryptedChkstr,
			"n": G.N,
			"type": G.VType,
			"os": "Windows 10",
			"name": "windows",
			"double_stack": "0",
			"_": "1602812428675",
		})
		ErrHandler(err)
		G.LoginResult,err=Util.GetResult(res)
		ErrHandler(err)
	}
	fmt.Println("The loggin result is: " + G.LoginResult)
}
