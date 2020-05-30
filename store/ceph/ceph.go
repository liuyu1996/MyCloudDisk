package ceph

import (
	"MyCloudDisk/config"
	"github.com/go-amz/amz/aws"
	"github.com/go-amz/amz/s3"
)

var cephConn *s3.S3

func GetCephConn() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}
	//1.初始化ceph信息
	auth := aws.Auth{
		AccessKey: config.Configs.AwsAccessKey,
		SecretKey: config.Configs.AwsSecretKey,
	}
	Region := aws.Region{
		Name:                 config.Configs.RegionName,
		EC2Endpoint:          config.Configs.EC2Endpoint,
		S3Endpoint:           config.Configs.S3Endpoint,
		S3BucketEndpoint:     config.Configs.S3BucketEndpoint,
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}
	//2.创建s3类型的连接
	return s3.New(auth, Region)
}

func GetCephBuket(bucket string) *s3.Bucket {
	conn := GetCephConn()
	return conn.Bucket(bucket)
}