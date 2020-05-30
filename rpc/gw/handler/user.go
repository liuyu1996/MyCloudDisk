package handler

import (
	"MyCloudDisk/config"
	"MyCloudDisk/rpc/proto"
	"MyCloudDisk/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"log"
	"net/http"
)

var (
	userCli proto.UserService
)

func init()  {
	service := micro.NewService()
	service.Init()

	//初始化一个rpcClient
	userCli = proto.NewUserService("go.micro.service.user", service.Client())
}

func SignUpHandler(c *gin.Context)  {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

func DoSignUpHandler(c *gin.Context)  {
	UserName := c.Request.FormValue("username")
	UserPwd := c.Request.FormValue("password")

	resp, err := userCli.Signup(context.TODO(), &proto.ReqSignup{
		Username: UserName,
		Password: UserPwd,
	})
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg" : resp.Message,
	})
}

func LoginHandler(c *gin.Context)  {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}


func DoSignInHandler(c *gin.Context)  {
	//1.校验用户名和密码
	userName := c.Request.FormValue("username")
	userPwd := c.Request.FormValue("password")
	hashPwd := utils.Sha1([]byte(userPwd+config.Configs.PwdSalt))

	resp, err := userCli.Signin(context.TODO(), &proto.ReqSignin{
		Username: userName,
		Password: hashPwd,
	})
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	rspMes := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token string
		}{
			Location: "/static/view/home.html",
			Username:userName,
			Token:resp.Token,
		},
	}
	c.JSON(http.StatusOK, rspMes)
}

func UserInfoHandler(c *gin.Context)  {
	//1.解析请求参数
	username := c.Request.FormValue("username")
	
	resp, err := userCli.UserInfo(context.TODO(), &proto.ReqUserInfo{
		Username: username,
	})
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	respMsg := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: gin.H{
			"Username": username,
			"SignupAt": resp.SignupAt,
			"LastActive" : resp.LastActiveAt,
		},
	}
	c.JSON(http.StatusOK, respMsg)
}

