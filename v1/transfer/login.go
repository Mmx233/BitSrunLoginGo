package srunTransfer

import "net/http"

type LoginForm struct {
	Domain   string `json:"domain"`
	UserName string `json:"username"`
	UserType string `json:"user_type"`
	PassWord string `json:"password"`
}

type LoginMeta struct {
	N    string `json:"n"`
	Type string `json:"type"`
	Acid string `json:"acid"`
	Enc  string `json:"enc"`
}

type LoginInfo struct {
	Form *LoginForm
	Meta *LoginMeta
}

type Login struct {
	//调用API时直接访问https URL
	Https bool
	//登录参数，不可缺省
	LoginInfo LoginInfo
	Client    *http.Client
}
