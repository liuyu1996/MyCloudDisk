package repository

import (
	"MyCloudDisk/db/mysql"
	"MyCloudDisk/models"
)

func UserFileUpload(userfile *models.UserFile) bool {
	if err := mysql.DBConn().Save(userfile).Error; err != nil{
		return false
	}
	return true
}

func GetUserFile(username string, limit int) ([]models.UserFile, error) {
	var userfileList []models.UserFile
	if err := mysql.DBConn().Where("user_name=?", username).
		Find(&userfileList).Limit(limit).Error; err != nil{
		return nil, err
	}
	return userfileList, nil
}
