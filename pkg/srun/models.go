package srun

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
	Form LoginForm
	Meta LoginMeta
}
