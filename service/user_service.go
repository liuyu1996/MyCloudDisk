package service

import (
	"MyCloudDisk/models"
	"MyCloudDisk/repository"
)

func AddUser (username string, pwd string) bool {
	user := models.User{
		UserName:       username,
		UserPwd:        pwd,
		Email:          "",
		Phone:          "",
		Status:         1,
	}
	return repository.Regist(user)
}

func GetUser (username string, hashPwd string) bool {
	user := repository.Login(username)
	if user == nil || user.UserPwd != hashPwd {
		return false
	}
	return true
}

func GetUserInfo(username string) (*models.User, error) {
	return repository.GetUserInfo(username)
}
