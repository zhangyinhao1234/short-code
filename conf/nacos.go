package conf

type Nacos struct {
	ServerAddr      string `yaml:"server-addr"`
	Port            uint64 `yaml:"server-port"`
	Namespace       string `yaml:"namespace"`
	DataId          string `yaml:"data-id"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	RefreshEnabled  bool   `yaml:"refresh-enabled"`
	RegisterEnabled bool   `yaml:"register-enabled"`
}
