package models

import "github.com/jinzhu/gorm"

type File struct {
	gorm.Model
	//Id       int64  `gorm:"id" json:"id"`
	FileSha1 string `gorm:"file_sha1" json:"file_sha1"` // 文件hash
	FileName string `gorm:"file_name" json:"file_name"` // 文件名
	FileSize int64  `gorm:"file_size" json:"file_size"` // 文件大小
	FileAddr string `gorm:"file_addr" json:"file_addr"` // 文件存储位置
	//CreateAt string `gorm:"create_at" json:"create_at"` // 创建日期
	//UpdateAt string `gorm:"update_at" json:"update_at"` // 更新日期
	Status   int64  `gorm:"status" json:"status"`       // 状态(可用/禁用/已删除等状态)
	Ext1     int64  `gorm:"ext1" json:"ext1"`           // 备用字段1
	Ext2     string `gorm:"ext2" json:"ext2"`           // 备用字段2
}

func (File) TableName() string {
	return "file"
}
