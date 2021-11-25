package srunTransfer

import "github.com/Mmx233/BitSrunLoginGo/models"

type LoginInfo struct {
	Form *srunModels.LoginForm
	Meta *srunModels.LoginMeta
}

type Login struct {
	//文件日志输出开关
	Demo bool
	//控制台日志打印开关
	OutPut bool
	//登陆前是否检查网络，只在离线时登录
	CheckNet bool
	//网络检查超时时间
	Timeout uint
	//登录参数，不可缺省
	LoginInfo LoginInfo
}
