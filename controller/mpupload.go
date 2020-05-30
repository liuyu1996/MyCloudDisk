package controller

import (
	mRedis "MyCloudDisk/cache/redis"
	"MyCloudDisk/models"
	"MyCloudDisk/mq"
	"MyCloudDisk/service"
	"MyCloudDisk/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type MultipartUploadInfo struct {
	FileHash string
	FileSize int
	UploadID string
	ChunkSize int
	ChunkCount int
}

//初始化分块上传
func InitialMultipartUploadHandler(c *gin.Context) {
	//1.解析用户请求参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":"Invalid request",
			"code": -1,
		})
		return
	}
	redisConn := mRedis.RedisClient.Pool.Get()
	if err := redisConn.Err(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":"Internal error",
			"code": -2,
		})
		return
	}
	defer redisConn.Close()
	//3.生成分块上传的初始化信息
	uploadInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:	filesize,
		UploadID:   username + fmt.Sprintf("%x",time.Now().UnixNano()),
		ChunkSize:  5*1024*1024,
		ChunkCount: int(math.Ceil(float64(filesize)/(5*1024*1024))),
	}
	//4.将初始化信息写入到redis缓存
	_, _ = redisConn.Do("HMSET", "MP_"+uploadInfo.UploadID,
		"chunkcount", uploadInfo.ChunkCount,
		"filehash", uploadInfo.FileHash,
		"filesize", uploadInfo.FileSize)
	//5.将响应初始化数据返回到客户端
	resp := utils.NewRespMsg(0, "OK", uploadInfo)
	c.JSON(http.StatusOK, resp)
}

func UploadPartHandler(c *gin.Context)  {
	//1.解析用户请求参数
	//username := r.Form.Get("username")
	uploadID := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")
	//2.获得redis连接
	redisConn := mRedis.RedisClient.Pool.Get()
	defer redisConn.Close()
	//3.获得文件句柄，用于存储分块内容
	filePath := "/data/" + uploadID + "/" + chunkIndex
	_ = os.MkdirAll(path.Dir(filePath), 0744)
	filestream, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "upload failed",
			"code" : -1,
		})
		return
	}
	defer filestream.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		if err != nil {
			break
		}
		_, err = filestream.Write(buf[:n])
		if err != nil {
			break
		}
	}
	//4.更新redis缓存状态
	_, err = redisConn.Do("HMSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)
	if err != nil {

	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "OK",
		"code" : 0,
	})
}
func CompleteUploadHandler(c *gin.Context)  {
	//1.解析请求参数
	uploadID := c.Request.FormValue("uploaadid")
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	filesize := c.Request.FormValue("filesize")
	//2。获得redis连接池中的一个连接
	redisConn := mRedis.RedisClient.Pool.Get()
	defer redisConn.Close()
	//3.通过uploadID查询redis并判断是否所有分块上传完成
	data, err := redis.Values(redisConn.Do("HGETALL", "MP_" + uploadID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "complete upload failed",
			"code" : -1,
		})
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data) ; i+=2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		}else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount ++
		}
	}
	if totalCount != chunkCount {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Invalid request",
			"code" : -1,
		})
		return
	}
	//4.合并分块
	tmpPath := "./data/" + uploadID
	fileAddr := "./tmp/"
	files, err := ioutil.ReadDir(tmpPath)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "合并分块失败",
			"code" : -1,
		})
		return
	}
	//创建文件，拿到文件句柄
	complateFile, err := os.Create(fileAddr + filename)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "complete upload failed",
			"code" : -1,
		})
		return
	}
	defer complateFile.Close()
	for _, f := range files{
		if f.Name() == ".DS_Store" {
			continue
		}
		//读取分块文件
		buffer, err := ioutil.ReadFile(fileAddr + f.Name())
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, gin.H{
				"msg": "complete upload failed",
				"code" : -1,
			})
			return
		}
		_, _ = complateFile.Write(buffer)
	}
	//4.将文件异步写入oss
	ossPath := "oss/" + filehash
	transData := mq.TransferData{
		FileHash:      filehash,
		CurLocation:   fileAddr + filename,
		DestLocation:  ossPath,
		DestStoreType: "OSS",
	}
	pubData, _ := json.Marshal(transData)
	rabbitmq := mq.NewRabbitMQ()
	ok := rabbitmq.Publish(pubData)
	if !ok {
		//TODO:加入重新发送消息逻辑
	}

	//5.更新唯一文件表及用户文件表
	fsize, _ := strconv.Atoi(filesize)
	file := models.File{
		FileSha1: filehash,
		FileName: filename,
		FileSize: int64(fsize),
		FileAddr: "",
	}
	service.UploadFile(file)
	service.UserFileUpload(username, filehash, filename, int64(fsize))
	//6.响应处理结果
	c.JSON(http.StatusOK, gin.H{
		"msg": "OK",
		"code" : 0,
	})
}


//TODO：取消分块上传
func CancelUploadPartHandler(c *gin.Context)  {
	//1.删除已存在的分块文件
	//2.删除redis缓存状态
	//3.更新mysql文件status
}  

//TODO：查看分块上传进度
func MultipartUploadProgressHandler(c *gin.Context)  {
	//检查分块上传状态是否有效
	//获取分块初始化信息
	//获取已上传的分块信息
}