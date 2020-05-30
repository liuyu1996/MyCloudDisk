package repository

import (
	"MyCloudDisk/db/mysql"
	"MyCloudDisk/models"
)

func UploadFile(file *models.File) bool {
	if err := mysql.DBConn().Save(file).Error; err != nil{
		return false
	}
	return true
}

func GetFileInfo(fileHash string) (*models.File, error) {
	file := &models.File{}
	if err := mysql.DBConn().Where("file_sha1=? AND status=?", fileHash, 1).First(file).Error ;err != nil{
		return nil, err
	}
	return file, nil
}

func UpdateFile(file *models.File) bool {
	if err := mysql.DBConn().Update(file).Error; err != nil {
		return false
	}
	return true
}


func DeleteFile(filehash string) error {
	if err := mysql.DBConn().Where("file_sha1", filehash).Delete(models.File{}).Error; err != nil {
		return err
	}
	return nil
}