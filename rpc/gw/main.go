package main

import (
	"MyCloudDisk/cache/redis"
	"MyCloudDisk/db/mysql"
	"MyCloudDisk/rpc/gw/route"
)

func init() {
	mysql.Default()
	redis.Default()
}

func main()  {
	r := route.Router()
	_ = r.Run(":8080")
}
