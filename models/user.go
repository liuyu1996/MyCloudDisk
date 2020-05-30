package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	//Id             int64  `gorm:"id" json:"id"`
	UserName       string `gorm:"user_name" json:"user_name"`             // 用户名
	UserPwd        string `gorm:"user_pwd" json:"user_pwd"`               // 用户encoded密码
	Email          string `gorm:"email" json:"email"`                     // 邮箱
	Phone          string `gorm:"phone" json:"phone"`                     // 手机号
	EmailValidated int64  `gorm:"email_validated" json:"email_validated"` // 邮箱是否已验证
	PhoneValidated int64  `gorm:"phone_validated" json:"phone_validated"` // 手机号是否已验证
	//SignupAt       string `gorm:"signup_at" json:"signup_at"`             // 注册日期
	//LastActive     string `gorm:"-" json:"last_active"`         // 最后活跃时间戳
	Profile        string `gorm:"profile" json:"profile"`                 // 用户属性
	Status         int64  `gorm:"status" json:"status"`                   // 账户状态(启用/禁用/锁定/标记删除等)
}

func (User) TableName() string {
	return "user"
}
