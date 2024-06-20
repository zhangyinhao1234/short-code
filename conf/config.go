package conf

type AppConf struct {
	Server     Server     `yaml:"server"`
	Redis      Redis      `yaml:"redis"`
	Logger     Logger     `yaml:"log"`
	ShortCode  ShortCode  `yaml:"shortCode"`
	ClickHouse ClickHouse `yaml:"clickHouse"`
}
