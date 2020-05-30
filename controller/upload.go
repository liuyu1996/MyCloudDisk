package controller

import (
	"MyCloudDisk/config"
	"MyCloudDisk/models"
	"MyCloudDisk/mq"
	"MyCloudDisk/service"
	"MyCloudDisk/store/oss"
	"MyCloudDisk/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func UploadHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/index.html")
}

func DoUploadHandler(c *gin.Context) {
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Failed to get data : %s", err)
		return
	}
	defer file.Close()

	fileMeta := models.File{
		FileName: head.Filename,
		FileAddr: "./tmp/" + head.Filename,
	}

	newFile, err := os.Create(fileMeta.FileAddr)
	if err != nil {
		log.Printf("create file failed : %s", err)
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		log.Printf("Failed to save file : %s", err)
		return
	}

	_, _ = newFile.Seek(0, 0)
	fileMeta.FileSha1 = utils.FileSha1(newFile)

	//将文件异步写入oss
	ossPath := "oss/" + fileMeta.FileSha1
	data := mq.TransferData{
		FileHash:      fileMeta.FileSha1,
		CurLocation:   fileMeta.FileAddr,
		DestLocation:  ossPath,
		DestStoreType: "OSS",
	}
	pubData, _ := json.Marshal(data)
	rabbitmq := mq.NewRabbitMQ()
	ok := rabbitmq.Publish(pubData)
	if !ok {
		//TODO:加入重新发送消息逻辑
	}

	//将文件写入ceph存储
	//_, _ = newFile.Seek(0, 0)
	//data, _ := ioutil.ReadAll(newFile)

	//bucket := ceph.GetCephBuket("userfile")
	//cephPath := "/ceph/" + fileMeta.FileSha1
	//_ = bucket.Put(cephPath, data, "octet-stream", s3.PublicRead)
	//fileMeta.FileAddr = cephPath

	service.UploadFile(fileMeta)
	//更新用户文件表
	username := c.Request.FormValue("username")
	ok = service.UserFileUpload(username, fileMeta.FileSha1, fileMeta.FileName,
		fileMeta.FileSize)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": "upload failed",
			"code" : -1,
		})
	} else {
		c.Redirect(http.StatusFound, "/static/view/home.html")
	}
	c.Redirect(http.StatusFound, "/file/upload/suc")
}

func UploadSucHandler(c *gin.Context) {
	_, _ = c.Writer.Write([]byte("upload success"))
}

func GetFileMetaHandler(c *gin.Context) {
	filehash := c.Request.Form["filehash"][0]
	file, err := service.GetFileInfo(filehash)
	data, err := json.Marshal(file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Internal error",
			"code": -1,
		})
		return
	}
	_, _ = c.Writer.Write(data)
}

func QueryFileInfoHandler(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")
	fileList, err := service.GetUserFileList(username, limit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Internal error",
			"code": -2,
		})
		return
	}
	data, err := json.Marshal(fileList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Internal error",
			"code": -2,
		})
		return
	}
	_, _ = c.Writer.Write(data)
}

func DownloadHandler(c *gin.Context) {
	filesha1 := c.Request.FormValue("filehash")
	filemeta, err := service.GetFileInfo(filesha1)

	f, err := os.Open(filemeta.Location)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Internal error",
			"code": -2,
		})
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Internal error",
			"code": -2,
		})
		return
	}
	c.Header("Content-Type", "application/octect-stream")
	c.Header("Content-Description", "attachment;filename=\""+filemeta.FileName+"\"")
	_, _ = c.Writer.Write(data)
}

func UpdateFileMetaHandler(c *gin.Context) {
	opType := c.Request.FormValue("op")
	filehash := c.Request.FormValue("filehash")
	fileName := c.Request.FormValue("filename")

	if opType != "0" {
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}
	if c.Request.Method != "POST" {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
	}

	fileInfo, err := service.GetFileInfo(filehash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg" : "Internal error",
			"code" : -2,
		})
		return
	}
	fileInfo.FileName = fileName
	service.UpdateFileInfo(fileInfo)

	data, err := json.Marshal(fileInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg" : "Internal error",
			"code" : -2,
		})
		return
	}
	c.JSON(http.StatusOK, data)
	//w.WriteHeader(http.StatusOK)
	//w.Write(data)
}

func DeleteFileHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	ok := service.DeleteFile(filehash)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg" : "Internal error",
			"code" : -2,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg" : "delete success",
		"code" : 0,
	})
}

func FastUploadHandler(c *gin.Context) {
	//1.解析请求参数
	username := c.Request.FormValue("username")
	fileHash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	//2.从文件表中查询相同hash的记录
	fileInfo, err := service.GetFileInfo(fileHash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg" : "Internal error",
			"code" : -2,
		})
		return
	}
	//3.查不到记录返回失败
	if fileInfo == nil {
		resp := utils.RespMsg{
			Code: -1,
			Msg:  "秒传失败",
			Data: nil,
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	//4.查询到记录则将文件信息写入用户文件表
	ok := service.UserFileUpload(username, fileHash, filename, int64(filesize))
	if ok {
		resp := utils.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
			Data: nil,
		}
		c.JSON(http.StatusOK, resp)
		return
	} else {
		resp := utils.RespMsg{
			Code: -2,
			Msg:  "秒传失败，请稍后重试",
			Data: nil,
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

//生成文件下载地址并返回给客户端
func DownloadURLHandler(c *gin.Context) {
	filehash :=c.Request.FormValue("filehash")
	fileInfo, _ := service.GetFileInfo(filehash)
	signedURL := oss.DownloadUrl(config.Configs.OssBucket, fileInfo.Location)
	_, _ = c.Writer.Write([]byte(signedURL))
}
