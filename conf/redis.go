package conf

type Redis struct {
	Addrs        []string `yaml:"addrs"`
	Password     string   `yaml:"password"`
	PoolSize     int      `yaml:"poolSize"`
	MinIdleConns int      `yaml:"minIdleConns"`
}
