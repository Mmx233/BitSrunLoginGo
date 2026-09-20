package srun

type LoginForm struct {
	Domain   string `json:"domain" yaml:"domain"`
	Username string `json:"username" yaml:"username"`
	//运营商类型
	UserType string `json:"user_type" yaml:"user_type"`
	Password string `json:"password" yaml:"password"`
}

type LoginFormOverride struct {
	Domain   *string `json:"domain,omitempty" yaml:"domain,omitempty"`
	Username *string `json:"username,omitempty" yaml:"username,omitempty"`
	UserType *string `json:"user_type,omitempty" yaml:"user_type,omitempty"`
	Password *string `json:"password,omitempty" yaml:"password,omitempty"`
}

func (o LoginFormOverride) Apply(form LoginForm) LoginForm {
	if o.Domain != nil {
		form.Domain = *o.Domain
	}
	if o.Username != nil {
		form.Username = *o.Username
	}
	if o.UserType != nil {
		form.UserType = *o.UserType
	}
	if o.Password != nil {
		form.Password = *o.Password
	}
	return form
}

type LoginMeta struct {
	N           string `json:"n" yaml:"n"`
	Type        string `json:"type" yaml:"type"`
	Acid        string `json:"acid" yaml:"acid"`
	Enc         string `json:"enc" yaml:"enc"`
	OS          string `json:"os" yaml:"os"`
	Name        string `json:"name" yaml:"name"`
	InfoPrefix  string `json:"info_prefix" yaml:"info_prefix"`
	DoubleStack bool   `json:"double_stack" yaml:"double_stack"`
}

type LoginInfo struct {
	Form LoginForm
	Meta LoginMeta
}
