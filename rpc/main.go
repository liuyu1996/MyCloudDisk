package main

import (
	"MyCloudDisk/cache/redis"
	"MyCloudDisk/db/mysql"
	"MyCloudDisk/rpc/handler"
	"MyCloudDisk/rpc/proto"
	"github.com/micro/go-micro"
	"log"
)

func init() {
	mysql.Default()
	redis.Default()
}

func main()  {
	//创建一个service
	service := micro.NewService(
			micro.Name("go.micro.service.user"))
	service.Init()

	_ = proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	err := service.Run()
	if err != nil {
		log.Println(err)
	}
}


