package BitSrun

func GenerateLoginInfo(Form *LoginForm, Meta *LoginMeta) *LoginContext {
	return &LoginContext{
		Meta: Meta,
		Form: &LoginForm{
			UserName: func() string {
				if Form.UserType == "" {
					return Form.UserName
				} else {
					return Form.UserName + "@" + Form.UserType
				}
			}(),
			PassWord: Form.PassWord,
		},
	}
}
