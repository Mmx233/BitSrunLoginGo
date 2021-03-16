package Modles

type LoginForm struct {
	Domain   string `json:"domain"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type LoginMeta struct {
	N     string `json:"n"`
	VType string `json:"v_type"`
	Acid  string `json:"acid"`
	Enc   string `json:"enc"`
}

type LoginInfo struct {
	UrlLoginPage       string
	UrlGetChallengeApi string
	UrlLoginApi        string

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
