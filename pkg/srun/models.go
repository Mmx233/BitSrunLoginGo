package srun

import (
	"net/http"
)

type LoginForm struct {
	Domain   string `json:"domain"`
	UserName string `json:"username"`
	//运营商类型
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

type Conf struct {
	//调用 API 时直接访问 https URL
	Https bool
	//登录参数，不可缺省
	LoginInfo LoginInfo
	Client    *http.Client

	api Api
}

func (a *Conf) initApi() {
	a.api.Init(a.Https, a.LoginInfo.Form.Domain, a.Client)
}
