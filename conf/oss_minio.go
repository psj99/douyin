package conf

type MinioOSS struct {
	Service         string `yaml:"service"`
	OssHost         string `yaml:"ossHost"`
	OssPort         string `yaml:"ossPort"`
	BucketName      string `yaml:"bucketName"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	Expiry          int    `yaml:"expiry"`
}
