package main

import (
	"MyCloudDisk/cache/redis"
	"MyCloudDisk/config"
	"MyCloudDisk/db/mysql"
	"MyCloudDisk/route"
	"log"
	"os"
)

func init() {
	mysql.Default()
	redis.Default()
}

func main()  {
	router := route.Router()
	err := router.Run("0.0.0.0:" + config.Configs.UploadServicePort)
	if err != nil {
		log.Println("start service error!!")
		os.Exit(0)
	}
}
