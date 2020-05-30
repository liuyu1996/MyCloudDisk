package route

import (
	"MyCloudDisk/middleware"
	"MyCloudDisk/rpc/gw/handler"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")

	router.GET("/user/signup", handler.SignUpHandler)
	router.POST("/user/signup", handler.DoSignUpHandler)

	router.Use(middleware.CheckAuth())
	{
		router.GET("/user/signin", handler.LoginHandler)
		router.POST("/user/signin", handler.DoSignInHandler)

		router.POST("/user/info", handler.UserInfoHandler)
	}


	return router
}
