package service

import (
	"MyCloudDisk/models"
	"MyCloudDisk/repository"
)

func UserFileUpload(username, fileHash, filename string, fileSize int64) bool {
	userfile := models.UserFile{
		UserName: username,
		FileSha1: fileHash,
		FileSize: fileSize,
		FileName: filename,
		Status:   1,
	}
	return repository.UserFileUpload(&userfile)
}

func GetUserFileList(username string, limit int) ([]models.UserFile, error) {
	return repository.GetUserFile(username, limit)
}
