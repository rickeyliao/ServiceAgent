package file

import "github.com/rickeyliao/ServiceAgent/translayer/control"

type FileHead struct {
	hashcode []byte
	fileName string
	fileSize int64
}

type FileTransReq struct {
	FileHead
	control.NbsControlHead
}


type FileTransResp struct {
	FileHead
	startpos int64
}










