package shadowsock

import (
	"sync"
	"sync/atomic"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/rickeyliao/ServiceAgent/common"
	"time"
)

type ServerStatistics struct {
	upBytes int64
	downBytes int64
}

var (
	ssInst *ServerStatistics
	ssInstLock sync.Mutex

	quit chan int
)

func GetSSInst() *ServerStatistics {
	if ssInst == nil{
		ssInstLock.Lock()
		defer ssInstLock.Unlock()

		if ssInst == nil{
			quit = make(chan int,0)
			ssInst = &ServerStatistics{}
			ssInst.Load()
		}

	}

	return ssInst
}

func (ss *ServerStatistics)IncUP(rBytes int64)  {
	atomic.AddInt64(&ss.upBytes,rBytes)
}

func (ss *ServerStatistics)IncDown(rBytes int64)  {
	atomic.AddInt64(&ss.downBytes,rBytes)
}

func (ss *ServerStatistics)Save()  {
	bjson,err:=json.Marshal(*ss)

	if err!=nil{
		return
	}

	tools.Save2File(bjson,common.GetSAConfig().GetSSStatFile())
}

func (ss *ServerStatistics)Load()  {
	d,err:=tools.OpenAndReadAll(common.GetSAConfig().GetSSStatFile())
	if err!=nil{
		return
	}

	ss1:=&ServerStatistics{}

	err=json.Unmarshal(d,ss1)
	if err!=nil{
		return
	}

	*ss = *ss1
}

func (ss *ServerStatistics)GetUPBytes() int64 {
	return ss.upBytes
}

func (ss *ServerStatistics)GetDownBytes() int64  {
	return ss.downBytes
}

func (ss *ServerStatistics)IntervalSave()  {
	var cnt int64
	for{
		cnt ++
		if cnt % 300 == 0{
			ss.Save()
		}
		select {
		case <-quit:
			break
		default:
			time.Sleep(time.Second)
		}
	}

}

func (ss *ServerStatistics)Quit()  {
	quit<-1
	ss.Save()
}
