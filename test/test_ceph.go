package main

import (
	"MyCloudDisk/store/ceph"
	"fmt"
	"github.com/go-amz/amz/s3"
)

func main()  {
	bucket := ceph.GetCephBuket("testBucket")
	//创建一个新的bucket
	err := bucket.PutBucket(s3.PublicRead)
	if err != nil {
		fmt.Println(err)
		return
	}
	//查询这个bucket下面指定条件的object keys
	res, err := bucket.List("", "", "", 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("keys:%+v\n", res)
	//新上传一个对象
	err = bucket.Put("/testupload/test.txt", []byte("test"), "octet-stream", s3.PublicRead)
	if err != nil {
		fmt.Println(err)
		return
	}
	//查询这个bucket下面指定条件的object keys
	res, err = bucket.List("", "", "", 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("keys:%+v\n", res)
}
