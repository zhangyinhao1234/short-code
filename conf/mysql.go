package conf

type MySQL struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	PW       string `yaml:"pw"`
	Database string `yaml:"database"`
}
