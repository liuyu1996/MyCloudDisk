package oss

import (
	"MyCloudDisk/config"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
)

var ossCli *oss.Client

func UploadToOss(fileName string, path string, bn string) bool {
	client := Client()
	bucket, err := client.Bucket(bn)
	if err != nil {
		log.Printf("Getting bucket failed: %s", err)
		return false
	}
	err = bucket.UploadFile(fileName, path, 500*1024, oss.Routines(5))
	if err != nil {
		log.Printf("uploading object failed :%s", err)
		return false
	}
	return true
}

func Client() *oss.Client  {
	if ossCli != nil {
		return ossCli
	}
	ossCli, err := oss.New(config.Configs.OssEndPoint,
		config.Configs.OssAccessKey, config.Configs.OssAccessSecret)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return ossCli
}

func Bucket(bucketName string) *oss.Bucket  {
	client := Client()
	if client != nil {
		bucket, err := client.Bucket(bucketName)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		return bucket
	}
	return nil
}

func DownloadUrl(bucketName, objectName string) string {
	signedUrl, err := Bucket(bucketName).SignURL(objectName, oss.HTTPGet, 3600)
	if err != nil {
		return ""
	}
	return signedUrl
}

//对指定的bucket设置生命周期规则
func SetLifecycleRule(bucketName string, prefix string, days int)  {
	rule := oss.BuildLifecycleRuleByDays("rule", prefix, true, days)
	rules := []oss.LifecycleRule{rule}

	_ = Client().SetBucketLifecycle(bucketName, rules)
}