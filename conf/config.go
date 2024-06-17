package conf

type AppConf struct {
	Server     Server     `yaml:"server"`
	Redis      Redis      `yaml:"redis"`
	Logger     Logger     `yaml:"log"`
	ShotCode   ShotCode   `yaml:"shotCode"`
	ClickHouse ClickHouse `yaml:"clickHouse"`
}
