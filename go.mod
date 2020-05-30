module MyCloudDisk

go 1.14

require (
	github.com/aliyun/aliyun-oss-go-sdk v2.1.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gohouse/converter v0.0.3
	github.com/golang/protobuf v1.4.1
	github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
	github.com/jinzhu/gorm v1.9.12
	github.com/micro/go-micro v1.18.0
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/protobuf v1.24.0
	gopkg.in/ini.v1 v1.57.0
)
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
