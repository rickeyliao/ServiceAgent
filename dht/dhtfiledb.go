package dht

import (
	"github.com/rickeyliao/ServiceAgent/db"
	"sync"
	"github.com/rickeyliao/ServiceAgent/common"
	"path"
	"github.com/pkg/errors"
	"github.com/kprc/nbsnetwork/tools"
	"encoding/json"
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
	IsLocal  bool		`json:"islocal"`
	LastAccessTime int64 `json:"time"`
}

func Insert(fhashv string,islocal bool) error {
	if fhashv == "" && !common.CheckNbsCotentHash(fhashv){
		return errors.New("Content Hash Error")
	}
	fsd:=&FileStoreDesc{IsLocal:islocal,LastAccessTime:tools.GetNowMsTime()}

	if bfsd,err:=json.Marshal(fsd);err!=nil{
		return err
	}else{
		GetFileStoreDB().Update(fhashv,string(bfsd))
	}

	return nil
}

func Update(fhashv string,islocal bool) error {
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


