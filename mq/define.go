package mq

//消息结构体
type TransferData struct {
	FileHash string
	CurLocation string
	DestLocation string
	DestStoreType string
}
