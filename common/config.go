package common

import (
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"sync"
	"path"
	"os"
	"encoding/json"
)


type SAConfig struct {
	DownloadDir string
	UploadDir string
	UploadMaxSize int64
	RemoteServerIP string
	RemoteServerPort uint16
	VerifyPath string
	ConsumePath string
	EmailPath string
	UpdateClientSoftwarePath string
	SoftWareVersion string
	BootstrapIPAddress []string
	TestIPAddress string
}


type SARootConfig struct {
	HomeDir     string
	CfgDir      string
	CfgFileName string
	SacInst *SAConfig
}

var (
	sarInst *SARootConfig
	saclock sync.Mutex
)

func GetSARootCfg() *SARootConfig {
	if sarInst == nil{
		saclock.Lock()
		defer saclock.Unlock()

		if sarInst == nil{
			sarInst = DefaultInitRootConfig()
		}

	}
	return sarInst
}

func DefaultInitRootConfig() *SARootConfig {
	usrdir,err := tools.Home()
	if err !=nil{
		log.Fatal("Get Home Dir failed")
	}

	homedir:= path.Join(usrdir,".sa")


	cfgdir:=path.Join(homedir,"config")

	//The config path looks like /home/user/.sa/config
	return &SARootConfig{HomeDir:homedir,CfgDir:cfgdir,CfgFileName:"sa.json"}
}

func DefaultInitConfig() *SAConfig  {
	sa:=&SAConfig{}

	sa.BootstrapIPAddress = []string{"103.45.98.72:50810","174.7.124.45:50810"}
	sa.ConsumePath = "/public/keys/consume"
	sa.DownloadDir = "download"
	sa.EmailPath = "/public/key/refresh"
	sa.RemoteServerIP = "47.90.242.83:80"
	sa.RemoteServerPort = 50810
	sa.SoftWareVersion = "0.1.0.0513"
	sa.UpdateClientSoftwarePath = "/public/app"
	sa.UploadDir = "upload"
	sa.UploadMaxSize = 1000    //1g
	sa.VerifyPath = "/public/keys/verify"
	sa.TestIPAddress = "/localipaddress"

	return sa
}


func (sar *SARootConfig)InitConfig() *SARootConfig  {
	if sar.HomeDir == "" || sar.CfgDir == "" || sar.CfgFileName == "" {
		log.Fatal("Please Set Config Path")
	}

	cfgname:=path.Join(sar.CfgDir,sar.CfgFileName)

	if !tools.FileExists(cfgname){
		if !tools.FileExists(sar.CfgDir){
			os.MkdirAll(sar.CfgDir,0755)
		}
		sac := DefaultInitConfig()
		sar.SacInst = sac
		bjson,err:=json.Marshal(*sac)
		if err!=nil{
			log.Fatal("Json module error")
		}
		tools.Save2File(bjson,cfgname)
    }else{
    	bjson,err:=tools.OpenAndReadAll(cfgname)
    	if err!=nil{
    		log.Fatal(err)
		}

    	sac := &SAConfig{}
    	err = json.Unmarshal(bjson,sac)
		if err!=nil{
			log.Fatal(err)
		}
    	sar.SacInst = sac
	}

	download:=path.Join(sar.HomeDir,sar.SacInst.DownloadDir)
	upload:=path.Join(sar.HomeDir,sar.SacInst.UploadDir)

	if !tools.FileExists(download){
		os.MkdirAll(download,0755)
	}
	if !tools.FileExists(upload){
		os.MkdirAll(upload,0755)
	}

	return sar
}

