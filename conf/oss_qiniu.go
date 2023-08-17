package conf

type QiNiuOSS struct {
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Bucket    string `yaml:"bucket"`
	Domain    string `yaml:"domain"`
}
