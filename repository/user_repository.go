package repository

import (
	"MyCloudDisk/db/mysql"
	"MyCloudDisk/models"
	"log"
)

func Regist(user models.User) bool {
	if err := mysql.DBConn().Create(&user).Error; err != nil{
		return false
	}
	return true
}

func Login(username string) *models.User {
	user := &models.User{}
	if err := mysql.DBConn().Where("user_name=?", username).First(user).Error; err != nil{
		log.Println(err.Error())
		return nil
	}
	return user
}


func GetUserInfo(username string) (*models.User, error) {
	user := &models.User{}
	if err := mysql.DBConn().Where("user_name=?", username).First(user).Error; err != nil{
		return nil, err
	}
	return user, nil
}