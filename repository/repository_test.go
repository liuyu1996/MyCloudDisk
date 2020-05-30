package repository

import (
	"MyCloudDisk/models"
	"testing"
)


func TestUserFileUpload(t *testing.T) {
	var userfile  models.UserFile
	userfile.FileName = "111"
	userfile.FileSize = 10
	userfile.FileSha1 = "123465"
	//userfile := models.UserFile{
	//	UserName: "123",
	//	FileSha1: "3745874452",
	//	FileSize: 10,
	//	FileName: "111",
	//	UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	//	Status:   1,
	//}
	ok := UserFileUpload(&userfile)
	if !ok {
		t.Error("failed")
	}
}

func TestUploadFile(t *testing.T) {
	filehash := "eqwsdfwfw"
	filename := "132"
	fileSize := 10
	fileAddr := "/temp/123"

	ok := UploadFile(filehash, filename, int64(fileSize), fileAddr)
	if !ok {
		t.Error("0")
	}

}