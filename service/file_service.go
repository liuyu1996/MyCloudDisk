package service

import (
	"MyCloudDisk/models"
	"MyCloudDisk/repository"
)

type FileInfo struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
}

func UploadFile(file models.File) bool {
	file.Status = 1
	return repository.UploadFile(&file)
}

func GetFileInfo(fileHash string) (*FileInfo, error) {
	file, err := repository.GetFileInfo(fileHash)
	if err != nil {
		return nil, err
	}
	fileInfo := FileInfo{
		FileSha1: file.FileSha1,
		FileName: file.FileName,
		FileSize: file.FileSize,
		Location: file.FileAddr,
	}
	return &fileInfo, nil
}

func UpdateFileInfo(fileInfo *FileInfo) bool {
	file := models.File{
		FileSha1: fileInfo.FileSha1,
		FileName: fileInfo.FileName,
		FileSize: fileInfo.FileSize,
		FileAddr: fileInfo.Location,
	}
	return repository.UpdateFile(&file)
}

func DeleteFile(filehash string) bool {
	if err := repository.DeleteFile(filehash); err != nil{
		return false
	}
	return true
}
