package conf

type Server struct {
	StartEnv string
	Port     uint64 `yaml:"port"`
	Name     string `yaml:"name"`
	GinMode  string `yaml:"ginMode"`
	Pprof    bool   `yaml:"pprof"`
}
