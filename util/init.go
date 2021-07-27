package util

import (
	"autoLogin/global"
	"autoLogin/models"
	"log"
	"os"
)

func init() {
	//配置文件初始化
	Path := "Config.json"
	var c models.Config
	if !File.Exists(Path) {
		if err := File.WriteJson(
			Path,
			c.FillDefault(),
		); err != nil {
			log.Println("创建配置文件失败:\n", err.Error())
			os.Exit(1)
		}
		log.Println("已生成配置文件，请编辑 'Config.json' 然后重试")
		os.Exit(0)
	}

	if err := File.ReadJson(Path, &c); err != nil {
		log.Println("读取配置文件失败:\n", err.Error())
		os.Exit(1)
	}

	_ = File.WriteJson(Path, c.FillDefault())

	global.Config = &c
}
