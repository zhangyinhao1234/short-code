package initialize

import (
	"embed"
	"gopkg.in/yaml.v3"
	"short-code/global"
)

//go:embed dev.yaml
//go:embed local.yaml

var f embed.FS

func ConfByEnv(env string) {
	if env == "" {
		panic("缺少 env 环境变量信息,无法加载配置文件")
	}
	global.CONF.Server.StartEnv = env
	data, _ := f.ReadFile(env + ".yaml")
	if err := yaml.Unmarshal(data, &global.CONF); err != nil {
		panic(err)
	}
}
