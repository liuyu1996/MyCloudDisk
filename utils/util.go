package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

func GetAppDir() string {
	appDir, err := os.Getwd()
	if err != nil {
		file, _ := exec.LookPath(os.Args[0])
		applicationPath, _ := filepath.Abs(file)
		appDir, _ = filepath.Split(applicationPath)
	}
	return appDir
}

func IsTimeStr(str string) bool {
	timeLayout := "2006-01-02 15:04:05"                        //转化所需模板
	loc, _ := time.LoadLocation("Local")                       //重要：获取时区
	theTime, err := time.ParseInLocation(timeLayout, str, loc) //使用模板在对应时区转化为time.time类型
	if err != nil {
		return false
	}
	if theTime.Unix() > 0 {
		return true
	}
	return false
}

//时间格式转换
func DateToDateTime(date string) string {
	timeTemplate := "2006-01-02T15:04:05+08:00" //常规类型
	toTemplate := "2006-01-02 15:04:05"
	stamp, _ := time.ParseInLocation(timeTemplate, date, time.Local)
	return time.Unix(stamp.Unix(), 0).Format(toTemplate)

}

//func SendMail(mailTo []string, subject string, body string) error {
//	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码
//
//	//mailConn := map[string]string{
//	//  "user": "xxx@163.com",
//	//  "pass": "your password",
//	//  "host": "smtp.163.com",
//	//  "port": "465",
//	//}
//
//	mailConn := map[string]string{
//		"user": "850838205@qq.com",
//		"pass": "hxaybppybgisbbae",
//		"host": "smtp.qq.com",
//		"port": "25",
//	}
//
//	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int
//
//	m := gomail.NewMessage()
//
//	m.SetHeader("From",  m.FormatAddress(mailConn["user"], "liuyu官方")) //这种方式可以添加别名，即“XX官方”
//	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
//	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
//	//m.SetHeader("From", mailConn["user"])
//	m.SetHeader("To", mailTo...)    //发送给多个用户
//	m.SetHeader("Subject", subject) //设置邮件主题
//	m.SetBody("text/html", body)    //设置邮件正文
//
//	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
//
//	err := d.DialAndSend(m)
//	return err
//
//}