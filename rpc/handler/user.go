package handler

import (
	"MyCloudDisk/config"
	"MyCloudDisk/rpc/proto"
	"MyCloudDisk/service"
	"MyCloudDisk/utils"
	"context"
	"net/http"
)

type User struct {
}

func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, resp *proto.RespSignup) error {
	username := req.Username
	password := req.Password

	//参数校验
	if len(username) < 3 || len(password) < 5 {
		resp.Code = -1
		resp.Message = "注册参数无效"
		return nil
	}
	hashPwd := utils.Sha1([]byte(password+config.Configs.PwdSalt))
	ok := service.AddUser(username, hashPwd)
	if ok {
		resp.Code = http.StatusOK
		resp.Message = "注册成功"
	}else {
		resp.Code = -2
		resp.Message = "注册失败"
	}
	return nil
}

func (u *User) Signin(ctx context.Context, req *proto.ReqSignin, resp *proto.RespSignin) error  {
	userName := req.Username
	userPwd := req.Password
	hashPwd := utils.Sha1([]byte(userPwd+config.Configs.PwdSalt))
	ok := service.GetUser(userName, hashPwd)
	if !ok {
		resp.Code = -1
		resp.Message = "用户未注册"
		resp.Token = ""
		return nil
	}
	token, err := utils.GenerateToken(userName)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "Internal error"
		resp.Token = ""
		return nil
	}else {
		resp.Code = http.StatusOK
		resp.Message = "登录成功"
		resp.Token = token
	}
	return nil
}

func (u *User) UserInfo(ctx context.Context,req *proto.ReqUserInfo, resp*proto.RespUserInfo) error  {
	//1.解析请求参数
	username := req.Username

	//3.查询用户信息
	user, err := service.GetUserInfo(username)
	if err != nil {
		resp.Code = -1
		resp.Message = "用户未注册"
		return nil
	}else {
		resp.Code = http.StatusOK
		resp.Message = "OK"
		resp.Username = user.UserName
		resp.Email = user.Email
		resp.Phone = user.Phone
		resp.SignupAt = user.CreatedAt.Format("2006-01-02 15:04:05")
		resp.LastActiveAt = user.UpdatedAt.Format("2006-01-02 15:04:05")
		resp.Status = int32(user.Status)
	}
	return nil
}
