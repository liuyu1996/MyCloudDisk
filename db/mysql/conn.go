package mysql

import (
	"MyCloudDisk/config"
	"MyCloudDisk/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
	err error
)


func Default() {
	db, err = gorm.Open("mysql",  config.Configs.DBUser + ":"+ config.Configs.DBPassword + "@tcp("+ config.Configs.DBHost + ":" + config.Configs.DBPort + ")/"+ config.Configs.DBName +"?charset=utf8&parseTime=true")
	//db, err = gorm.Open("mysql", "root:123456@tcp(www.secretbaseofly.cn:3308)/CloudDisk?charset=utf8")
	if err != nil {
		panic(err.Error() + config.Configs.DBDriver)
	}
	db.SingularTable(true)
	db.AutoMigrate(&models.User{},&models.UserFile{}, &models.File{})
}

func DBConn() *gorm.DB {
	return db
}