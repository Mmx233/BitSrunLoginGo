package srun

type LoginForm struct {
	Domain   string `json:"domain" yaml:"domain"`
	Username string `json:"username" yaml:"username"`
	//运营商类型
	UserType string `json:"user_type" yaml:"userType"`
	Password string `json:"password" yaml:"password"`
}

type LoginMeta struct {
	N    string `json:"n" yaml:"n"`
	Type string `json:"type" yaml:"type"`
	Acid string `json:"acid" yaml:"acid"`
	Enc  string `json:"enc" yaml:"enc"`
}

type LoginInfo struct {
	Form LoginForm
	Meta LoginMeta
}
