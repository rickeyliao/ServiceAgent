package common

import (
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"sync"
	"path"
	"os"
	"encoding/json"
	"github.com/spf13/viper"
	"strings"
)


type SAConfig struct {
	DownloadDir string
	UploadDir string
	UploadMaxSize int64
	RemoteServerIP string
	RemoteServerPort uint16
	VerifyPath string
	ConsumePath string
	PostSocks5Path string
	ListIpsPath string
	EmailPath string
	UpdateClientSoftwarePath string
	SoftWareVersion string
	LocalListenPort uint16
	BootstrapIPAddress []string
	TestIPAddress string
	ListenTyp string

}


type SARootConfig struct {
	HomeDir     string
	CfgDir      string
	CfgFileName string
	SacInst *SAConfig
}


type Roothome struct{
	Rootdir string
}

var (
	sarInst *SARootConfig
	saclock sync.Mutex
)

func GetSARootCfgHdir(hdir string,force bool) *SARootConfig {
	if sarInst == nil{
		saclock.Lock()
		defer saclock.Unlock()

		if sarInst == nil{
			sarInst = DefaultInitRootConfig(hdir,force)
		}

	}
	return sarInst
}

func GetSARootCfg() *SARootConfig  {
	if sarInst == nil{
		saclock.Lock()
		defer saclock.Unlock()

		if sarInst == nil{
			sarInst = DefaultInitRootConfig("",false)
		}

	}
	return sarInst
}

func forceInitRootConfig(hdir string) *SARootConfig  {
	var sahome string
	var homedir string
	var err error

	if homedir, err = tools.Home(); err != nil {
		log.Fatal("Can't Get Home Directory")
	}

	if hdir == ""{
		viper.AutomaticEnv()
		sahome=viper.GetString("sahome")
		if sahome == "" {
			sahome = path.Join(homedir,".sa")
		}
	}else{
		hdir=path.Clean(hdir)
		if isroot:=path.IsAbs(hdir);!isroot{
			sahome = path.Join(homedir,hdir)
		}else{
			sahome = hdir
		}
	}

	if sahome == "" {
		return nil
	}

	//save homedir to .sainit file
	rh:=Roothome{sahome}
	var brh []byte
	if brh,err=json.MarshalIndent(rh,"","\t");err!=nil{
		log.Fatal("Can't save to .sainit file")
	}
	tools.Save2File(brh,path.Join(homedir,".sainit"))

	cfgdir := path.Join(sahome,"config")

	return &SARootConfig{HomeDir:sahome,CfgDir:cfgdir,CfgFileName:"sa.json"}
}

func unforceInitRootConfig(hdir string) *SARootConfig {
	var homedir string
	var savedir string
	var d []byte
	var err error

	if homedir, err = tools.Home(); err != nil {
		log.Fatal("Can't Get Home Directory")
	}
	d,err =tools.OpenAndReadAll(path.Join(homedir,".sainit"))
	if d == nil || len(d) == 0{

		if hdir == ""{
			viper.AutomaticEnv()
			savedir=viper.GetString("sahome")
			if savedir=="" {
				savedir = path.Join(homedir, ".sa")
			}
		}else{
			if isroot:=path.IsAbs(hdir);!isroot{
				savedir = path.Join(homedir,hdir)
			}else{
				savedir = hdir
			}
		}
		//save homedir to .sainit file
		rh:=Roothome{savedir}
		var brh []byte
		if brh,err=json.MarshalIndent(rh,"","\t");err!=nil{
			log.Fatal("Can't save to .sainit file")
		}

		tools.Save2File(brh,path.Join(homedir,".sainit"))

	}else{
		prh:=&Roothome{}
		if err=json.Unmarshal(d,prh);err!=nil{
			log.Fatal("Cant recover home dir")
		}

		savedir = prh.Rootdir

	}
	if savedir == ""{
		return nil
	}
	return &SARootConfig{HomeDir:savedir,CfgDir:path.Join(savedir,"config"),CfgFileName:"sa.json"}
}


func DefaultInitRootConfig(hdir string,force bool) *SARootConfig {
	var sar *SARootConfig
	if force{
		sar =  forceInitRootConfig(hdir)
	}else {
		sar =  unforceInitRootConfig(hdir)
	}

	log.Println("Config Root:",sar.HomeDir)

	return sar
}

func DefaultInitConfig() *SAConfig  {
	sa:=&SAConfig{}

	sa.BootstrapIPAddress = []string{"103.45.98.72:50810","174.7.124.45:50810"}
	sa.ConsumePath = "/public/keys/consume"
	sa.DownloadDir = "download"
	sa.EmailPath = "/public/key/refresh"
	sa.RemoteServerIP = "207.148.9.49"
	sa.RemoteServerPort = 80
	sa.LocalListenPort = 50810
	sa.ListenTyp = "udp4"
	sa.SoftWareVersion = "0.1.0.0521"
	sa.UpdateClientSoftwarePath = "/public/app"
	sa.UploadDir = "upload"
	sa.UploadMaxSize = 1000    //1g
	sa.VerifyPath = "/public/keys/verify"
	sa.TestIPAddress = "/localipaddress"
	sa.ListIpsPath = "/public/servers/list"
	sa.PostSocks5Path = "/postsocks5"

	return sa
}

func (sar *SARootConfig)LoadCfg() *SAConfig  {
	viper.AddConfigPath(path.Join(sar.CfgDir))
	strarr:=strings.Split(sar.CfgFileName,".")
	viper.SetConfigName(strarr[0])

	if err:=viper.ReadInConfig();err!=nil{
		log.Println("Read config file error")
		os.Exit(1)
	}

	cfg:=&SAConfig{}
	viper.Unmarshal(cfg)

	sar.SacInst = cfg

	return cfg
}


func (sar *SARootConfig)IsInitialized() bool  {
	cfgname := path.Join(sar.CfgDir,sar.CfgFileName)
	if !tools.FileExists(cfgname){
		return false
	}

	return true
}



func (sar *SARootConfig)InitConfig(force bool) *SARootConfig  {
	if sar.HomeDir == "" || sar.CfgDir == "" || sar.CfgFileName == "" {
		log.Fatal("Please Set Config Path")
	}

	if force{
		os.RemoveAll(sar.HomeDir)
	}

	cfgname:=path.Join(sar.CfgDir,sar.CfgFileName)

	if !tools.FileExists(cfgname){
		if !tools.FileExists(sar.CfgDir){
			os.MkdirAll(sar.CfgDir,0755)
		}
		sac := DefaultInitConfig()
		sar.SacInst = sac
		bjson,err:=json.MarshalIndent(*sac,"","\t")
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

