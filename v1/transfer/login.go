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
	//文件日志输出开关
	Demo bool
	//控制台日志打印开关
	OutPut bool
	//登陆前是否检查网络，只在离线时登录
	CheckNet bool
	//登录参数，不可缺省
	LoginInfo LoginInfo
	Transport *http.Transport
}
