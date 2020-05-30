package controller

import (
	"MyCloudDisk/service"
	"MyCloudDisk/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type SignUpParams struct {
	UserName       string `gorm:"user_name" json:"username"`             // 用户名
	UserPwd        string `gorm:"user_pwd" json:"password"`               // 用户encoded密码
	//Email          string `gorm:"email" json:"email"`                     // 邮箱
	//Phone          string `gorm:"phone" json:"phone"`                     // 手机号
}

const(
	pwd_salt = "*#@850"
)

func SignUpHandler(c *gin.Context)  {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

func DoSignUpHandler(c *gin.Context)  {
	var params SignUpParams
	params.UserName = c.Request.FormValue("username")
	params.UserPwd = c.Request.FormValue("password")

	if len(params.UserName) < 3 || len(params.UserPwd) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":"Invalid parameter",
			"code": -1,
		})
		return
	}
	params.UserPwd  = utils.Sha1([]byte(params.UserPwd+pwd_salt))
	ok := service.AddUser(params.UserName, params.UserPwd)
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"msg":"SUCCESS",
			"code": 0,
		})
	}else {
		c.JSON(http.StatusOK, gin.H{
			"msg":"FAILED",
			"code": -2,
		})
	}
}

func LoginHandler(c *gin.Context)  {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

func DoSignInHandler(c *gin.Context)  {
	//1.校验用户名和密码
	userName := c.Request.FormValue("username")
	userPwd := c.Request.FormValue("password")
	hashPwd := utils.Sha1([]byte(userPwd+pwd_salt))
	ok := service.GetUser(userName, hashPwd)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg":"Login failed",
			"code": -1,
		})
		return
	}
	//2.生成访问凭证
	token, err := utils.GenerateToken(userName)
	log.Println(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":"Login failed",
			"code": -1,
		})
		return
	}
	//3.登录成功后重定向到首页
	resp := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token string
		}{
			Location: "/static/view/home.html",
			Username:userName,
			Token:token,
		},
	}
	c.JSON(http.StatusOK, resp)
}

func UserInfoHandler(c *gin.Context)  {
	//1.解析请求参数
	username := c.Query("username")

	//3.查询用户信息
	user, err := service.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":"Internal error",
			"code": -1,
		})
		return
	}
	resp := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	c.JSON(http.StatusOK, resp)
}

