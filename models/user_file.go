package models

import "github.com/jinzhu/gorm"

type UserFile struct {
	gorm.Model
	//Id         int64  `gorm:"id;primary_key" json:"id"`
	UserName   string `gorm:"user_name" json:"user_name"`
	FileSha1   string `gorm:"file_sha1" json:"file_sha1"`     // 文件hash
	FileSize   int64  `gorm:"file_size" json:"file_size"`     // 文件大小
	FileName   string `gorm:"file_name" json:"file_name"`     // 文件名
	//UploadAt   string `gorm:"upload_at" json:"upload_at"`     // 上传时间
	//LastUpdate string `gorm:"last_update" json:"last_update"` // 最后修改时间
	Status     int64  `gorm:"status" json:"status"`           // 文件状态(0正常1已删除2禁用)
}

func (UserFile) TableName() string {
	return "user_file"
}
