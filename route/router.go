package route

import (
	"MyCloudDisk/controller"
	"MyCloudDisk/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")

	//用户相关接口
	router.GET("/user/signup", controller.SignUpHandler)
	router.POST("/user/signup", controller.DoSignUpHandler)
	router.GET("/user/signin", controller.LoginHandler)
	router.POST("/user/signin", controller.DoSignInHandler)

	router.Use(middleware.CheckAuth())
	{
		router.POST("/user/info", controller.UserInfoHandler)
		//文件相关接口
		router.GET("/file/upload", controller.UploadHandler)
		router.POST("/file/upload", controller.DoUploadHandler)
		router.GET("/file/upload/suc", controller.UploadSucHandler)
		router.POST("/file/info", controller.GetFileMetaHandler)
		router.POST("/file/query", controller.QueryFileInfoHandler)
		router.POST("/file/download", controller.DownloadHandler)
		router.POST("/file/delete", controller.DeleteFileHandler)
		router.POST("/file/downloadurl", controller.DownloadURLHandler)

		router.POST("/file/fastupload", controller.FastUploadHandler)

		router.POST("/file/mpupload/init", controller.InitialMultipartUploadHandler)
		router.POST("/file/mpupload/uppart", controller.UploadPartHandler)
		router.POST("/file/mpupload/complete", controller.CompleteUploadHandler)

	}

	return router
}
