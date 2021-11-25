package srunModels

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
	UrlLoginPage       string
	UrlGetChallengeApi string
	UrlLoginApi        string
	UrlCheckApi        string

	Ip              string
	Token           string
	EncryptedInfo   string
	Md5             string
	EncryptedMd5    string
	EncryptedChkstr string
	LoginResult     string

	Form *LoginForm
	Meta *LoginMeta
}
