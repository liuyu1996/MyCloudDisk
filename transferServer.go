package main

import (
	"MyCloudDisk/config"
	"MyCloudDisk/mq"
	"MyCloudDisk/service"
	m_oss "MyCloudDisk/store/oss"
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"os"
)

func ProcessTransfer(msg[]byte) bool {
	//1.解析msg
	data := mq.TransferData{}
	err := json.Unmarshal(msg, data)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	//2..使用oss断点续传
	err = m_oss.Bucket(config.Configs.OssBucket).UploadFile(data.DestLocation,
		data.CurLocation,
		500*1024, oss.Routines(5),
		oss.Checkpoint(true, ""))
	//err = m_oss.Bucket(config.Configs.OssBucket).PutObject(
	//	data.DestLocation,
	//	bufio.NewReader(fileStream))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	//4.更新文件表的存储路径字段
	fileInfo := service.FileInfo{
		Location: data.DestLocation,
	}
	ok := service.UpdateFileInfo(&fileInfo)
	if !ok {
		return false
	}
	_ = os.Remove(data.CurLocation)
	return true
}


func main()  {
	log.Println("开始监听转移任务队列......")
	rabbitmq := mq.NewRabbitMQ()
	rabbitmq.StartConsume("", ProcessTransfer)
}