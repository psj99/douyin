package conf

type System struct {
	AppEnv   string `yaml:"appEnv"`
	HttpPort string `yaml:"httpPort"`
	TempDir  string `yaml:"tempDir"`
}
