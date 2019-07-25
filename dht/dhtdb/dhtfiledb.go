package dhtdb

import (
	"github.com/rickeyliao/ServiceAgent/db"
	"sync"
	"github.com/rickeyliao/ServiceAgent/common"
	"path"
	"github.com/pkg/errors"
	"github.com/kprc/nbsnetwork/tools"
	"encoding/json"
	"fmt"
)

var(
	filestoredb db.NbsDbInter
	filestoredblock sync.Mutex

	quit chan int
	wg *sync.WaitGroup
)

func GetFileStoreDB() db.NbsDbInter  {
	if filestoredb == nil{
		filestoredblock.Lock()
		defer filestoredblock.Unlock()
		if filestoredb == nil{
			filestoredb = newFileStoreDB()
		}
	}

	return filestoredb
}

func newFileStoreDB() db.NbsDbInter  {
	quit = make(chan int, 0)
	wg = &sync.WaitGroup{}

	cfg:=common.GetSAConfig()

	return db.NewFileDb(path.Join(cfg.GetFileDbDir(),cfg.FileStoreDB)).Load()
}


type FileStoreDesc struct {
	FileHash string		`json:"-"`
	IsLocal  bool		`json:"l"`
	LastAccessTime int64 `json:"t"`
	IsShare bool		`json:"s"`
	NodeAddr []string   `json:"na"`
	FileExists bool     `json:"e"`
}

func Find(fhashv string) (v string, err error)  {
	return  GetFileStoreDB().Find(fhashv)
}

func Insert(fhashv string,islocal bool,share bool,naddrs []string,exist bool) error {
	if fhashv == "" && !common.CheckNbsCotentHash(fhashv){
		return errors.New("Content Hash Error")
	}
	fsd:=&FileStoreDesc{IsLocal:islocal,LastAccessTime:tools.GetNowMsTime()}
	fsd.FileExists = exist
	fsd.IsShare = share
	if len(naddrs)>0{
		fsd.NodeAddr = naddrs
	}else{
		fsd.NodeAddr = make([]string,0)
	}

	if bfsd,err:=json.Marshal(fsd);err!=nil{
		return err
	}else{
		GetFileStoreDB().Update(fhashv,string(bfsd))
	}

	return nil
}

func UpdateTime(fhashv string) error {
	sfsd,err:=GetFileStoreDB().Find(fhashv)
	if err!=nil{
		return err
	}
	fsd:=&FileStoreDesc{}
	if err=json.Unmarshal([]byte(sfsd),fsd);err!=nil{
		return err
	}
	fsd.LastAccessTime = tools.GetNowMsTime()

	var bfsd []byte
	bfsd,err=json.Marshal(fsd)
	if err!=nil{
		return err
	}
	GetFileStoreDB().Update(fhashv,string(bfsd))
	return nil
}

func UpdateLocal(fhashv string,islocal bool) error {
	sfsd,err:=GetFileStoreDB().Find(fhashv)
	if err!=nil{
		return err
	}
	fsd:=&FileStoreDesc{}
	if err=json.Unmarshal([]byte(sfsd),fsd);err!=nil{
		return err
	}

	fsd.LastAccessTime = tools.GetNowMsTime()
	fsd.IsLocal = islocal

	var bfsd []byte
	bfsd,err=json.Marshal(fsd)
	if err!=nil{
		return err
	}
	GetFileStoreDB().Update(fhashv,string(bfsd))
	return nil
}

func Update(fhashv string,islocal,isshare bool, nbsaddrs []string,exist bool) error  {
	sfsd,err:=GetFileStoreDB().Find(fhashv)
	if err!=nil{
		return err
	}
	fsd:=&FileStoreDesc{}
	if err=json.Unmarshal([]byte(sfsd),fsd);err!=nil{
		return err
	}
	fsd.LastAccessTime = tools.GetNowMsTime()
	fsd.IsLocal = islocal
	fsd.IsShare = isshare
	fsd.FileExists = exist
	if fsd.NodeAddr == nil{
		fsd.NodeAddr = nbsaddrs
	}else{
		for _,addr:=range nbsaddrs{
			fsd.NodeAddr = append(fsd.NodeAddr,addr)
		}
	}

	var bfsd []byte
	bfsd,err=json.Marshal(fsd)
	if err!=nil{
		return err
	}
	GetFileStoreDB().Update(fhashv,string(bfsd))

	return nil
}

func GetFileExistFlag(key string) bool {
	if v,err:=Find(key);err==nil{
		return false
	}else{
		fsd:=&FileStoreDesc{}
		if err=json.Unmarshal([]byte(v),fsd);err!=nil{
			return false
		}

		return fsd.FileExists
	}

}


func DeleteAddr(fhashv,nbsaddr string) error {
	sfsd,err:=GetFileStoreDB().Find(fhashv)
	if err!=nil{
		return err
	}
	fsd:=&FileStoreDesc{}
	if err=json.Unmarshal([]byte(sfsd),fsd);err!=nil{
		return err
	}
	fsd.LastAccessTime = tools.GetNowMsTime()

	if len(fsd.NodeAddr) > 0{
		idx := len(fsd.NodeAddr)
		for i, addr := range fsd.NodeAddr {
			if addr == nbsaddr {
				idx = i
				break
			}
		}
		if idx<len(fsd.NodeAddr){
			fsd.NodeAddr = append(fsd.NodeAddr[:len(fsd.NodeAddr)],fsd.NodeAddr[len(fsd.NodeAddr)+1:]...)
		}
	}

	var bfsd []byte
	bfsd,err=json.Marshal(fsd)
	if err!=nil{
		return err
	}
	GetFileStoreDB().Update(fhashv,string(bfsd))

	return nil

}

func Print()  {
	it:=GetFileStoreDB().DBIterator()

	for{
		k,v:=it.Next()
		if k==""{
			break
		}

		fmt.Println(k,v)
	}

}

func Remove(key string)  {
	GetFileStoreDB().Delete(key)
}

func SaveDb()  {
	GetFileStoreDB().Save()
}