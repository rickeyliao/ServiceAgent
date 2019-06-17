package dht

import "sync"

type msgcnt struct {
	cnt uint64
}

var (
	msgcnt_inst *msgcnt
	msgcnt_lock sync.Mutex
)

func GetNextMsgCnt() uint64 {
	if msgcnt_inst == nil{
		msgcnt_lock.Lock()
		if msgcnt_inst == nil{
			msgcnt_inst = &msgcnt{cnt:0x20151031}
		}
		msgcnt_lock.Unlock()
	}

	cnt:=msgcnt_inst.cnt
	msgcnt_inst.cnt ++

	return cnt
}

