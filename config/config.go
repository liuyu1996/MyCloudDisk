package config

import (
	"gopkg.in/ini.v1"
	"os"
	"os/exec"
	"path/filepath"
)

type SysConfig struct {
	Env        string `ini:"env"`
	Debug      bool   `ini:"debug"`
	DBDriver   string `ini:"db_driver"`
	DBHost     string `ini:"db_host"`
	DBPort     string `ini:"db_port"`
	DBUser     string `ini:"db_user"`
	DBPassword string `ini:"db_password"`
	DBName     string `ini:"db_name"`

	ModelPath   string `ini:"model_path"`
	Tag         string `ini:"tag"`
	TablePrefix string `ini:"table_prefix"`
	Table       string `ini:"table"`

	RedisHost string `ini:"redis_host"`
	RedisPwd  string `ini:"redis_password"`
	//RedisDb           int    `ini:"redis_db"`
	//RedisCacheVersion string `ini:"redis_cache_version"`

	OssBucket       string `ini:"oss_bucket"`
	OssEndPoint     string `ini:"oss_endpoint"`
	OssAccessKey    string `ini:"oss_accesskey"`
	OssAccessSecret string `ini:"oss_accesskeysecret"`

	AsyTransferEnable int `ini:"asy_transfer_enable"`
	RabbitURL string	`ini:"rabbit_url"`
	TransExchangeName string `ini:"trans_exchangename"`
	TransOSSQueueName string `ini:"trans_ossqueueName"`
	TransOSSErrQueueName string `ini:"trans_osserrqueueName"`
	TransOSSRoutingKey string  `ini:"trans_ossroutingkey"`

	RegionName       string `ini:"region_name"`
	AwsAccessKey     string `ini:"aws_accessKey"`
	AwsSecretKey     string `ini:"aws_secretKey"`
	EC2Endpoint      string `ini:"EC2Endpoint"`
	S3Endpoint       string `ini:"S3Endpoint"`
	S3BucketEndpoint string `ini:"S3BucketEndpoint"`

	UploadServicePort string `ini:"uploadService_port"`

	JwtSecret string `ini:"jwt_secret"`

	PwdSalt string `ini:"pwd_salt"`
}

var Configs *SysConfig = &SysConfig{}

//加载系统配置文件
func init() {
	config := &SysConfig{}
	conf, err := ini.Load(getAppDir() + "/config.ini") //加载配置文件
	if err != nil {
		panic(err)
	}
	conf.BlockMode = false
	err = conf.MapTo(&config) //解析成结构体
	if err != nil {
		panic(err)
	}
	Configs = config
}

func getAppDir() string {
	appDir, err := os.Getwd()
	if err != nil {
		file, _ := exec.LookPath(os.Args[0])
		applicationPath, _ := filepath.Abs(file)
		appDir, _ = filepath.Split(applicationPath)
	}
	return appDir
}
