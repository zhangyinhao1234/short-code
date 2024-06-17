package conf

type Logger struct {
	FileDir  string `yaml:"fileDir"`
	Level    string `yaml:"level"`
	LinkName string `yaml:"linkName"`
}
