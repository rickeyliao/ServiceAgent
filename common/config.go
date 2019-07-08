package common

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/tools/crypt/nbscrypt"
	"github.com/mr-tron/base58"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"sync"
	"github.com/rickeyliao/ServiceAgent/htmlfile"
)

type SAConfig struct {
	DownloadPath             string          `json:"downloadpath"`
	Uploadpath               string          `json:"uploadpath"`
	TestIPAddressPath        string          `json:"testipaddress"`
	PostSocks5Path           string          `json:"postsocks5path"`
	VerifyPath               string          `json:"verifypath"`
	ConsumePath              string          `json:"consumepath"`
	ListIpsPath              string          `json:"listipspath"`
	EmailPath                string          `json:"emailpath"`
	UpdateClientSoftwarePath string          `json:"updateclientsoftwarepath"`
	KeyDir                   string          `json:"keydir"`
	PidDir                   string          `json:"piddir"`
	FileDBDir                string          `json:"filedbdir"`
	HomeIPDBFile             string          `json:"homeipdbfile"`
	FileStoreDir             string          `json:"filestoredir"`
	RemoteServerIP           string          `json:"remoteserverip"`
	RemoteServerPort         uint16          `json:"remoteserverport"`
	HttpListenPort           uint16          `json:"httplistenport"`		//50810 tcp/http for file transfer
	BootstrapIPAddress       []string        `json:"bootstrapipaddress"`
	ListenTyp                string          `json:"listentyp"`
	NbsRsaAddr               string          `json:"nbsaddr"`
	CmdListenIP              string          `json:"cmdlistenip"`
	CmdListenPort            uint16          `json:"cmdlistenport"`   //50811 tcp for cmd
	DhtListenPort            uint16          `json:"dhtlistenport"`	  //50811 udp for control message
	StaticFileDir            string			 `json:"staticfiledir"`
	LoginDir                 string          `json:"logindir"`
	Loginfile                string          `json:"loginfile"`
	LicenseAdminUser         [][]string      `json:"licenseadminuser"`
	Role                     int64           `json:"-"`
	PrivKey                  *rsa.PrivateKey `json:"-"`
	Root					 *SARootConfig   `json:"-"`
}

type SARootConfig struct {
	HomeDir     string
	CfgDir      string
	CfgFileName string
	SacInst     *SAConfig
}

type Roothome struct {
	Rootdir string
}

var (
	sarInst *SARootConfig
	saclock sync.Mutex
)

func GetSARootCfgHdir(hdir string, force bool) *SARootConfig {
	if sarInst == nil {
		saclock.Lock()
		defer saclock.Unlock()

		if sarInst == nil {
			sarInst = DefaultInitRootConfig(hdir, force)
		}

	}
	return sarInst
}

func GetSARootCfg() *SARootConfig {
	if sarInst == nil {
		saclock.Lock()
		defer saclock.Unlock()

		if sarInst == nil {
			sarInst = DefaultInitRootConfig("", false)
		}

	}
	return sarInst
}

func GetSAConfig() *SAConfig {
	return GetSARootCfg().SacInst
}

func forceInitRootConfig(hdir string) *SARootConfig {
	var sahome string
	var homedir string
	var err error

	if homedir, err = tools.Home(); err != nil {
		log.Fatal("Can't Get Home Directory")
	}

	if hdir == "" {
		viper.AutomaticEnv()
		sahome = viper.GetString("sahome")
		if sahome == "" {
			sahome = path.Join(homedir, ".sa")
		}
	} else {
		hdir = path.Clean(hdir)
		if isroot := path.IsAbs(hdir); !isroot {
			sahome = path.Join(homedir, hdir)
		} else {
			sahome = hdir
		}
	}

	if sahome == "" {
		return nil
	}

	//save homedir to .sainit file
	rh := Roothome{sahome}
	var brh []byte
	if brh, err = json.MarshalIndent(rh, "", "\t"); err != nil {
		log.Fatal("Can't save to .sainit file")
	}
	tools.Save2File(brh, path.Join(homedir, ".sainit"))

	cfgdir := path.Join(sahome, "config")

	return &SARootConfig{HomeDir: sahome, CfgDir: cfgdir, CfgFileName: "sa.json"}
}

func unforceInitRootConfig(hdir string) *SARootConfig {
	var homedir string
	var savedir string
	var d []byte
	var err error

	if homedir, err = tools.Home(); err != nil {
		log.Fatal("Can't Get Home Directory")
	}
	d, err = tools.OpenAndReadAll(path.Join(homedir, ".sainit"))
	if d == nil || len(d) == 0 {

		if hdir == "" {
			viper.AutomaticEnv()
			savedir = viper.GetString("sahome")
			if savedir == "" {
				savedir = path.Join(homedir, ".sa")
			}
		} else {
			if isroot := path.IsAbs(hdir); !isroot {
				savedir = path.Join(homedir, hdir)
			} else {
				savedir = hdir
			}
		}
		//save homedir to .sainit file
		rh := Roothome{savedir}
		var brh []byte
		if brh, err = json.MarshalIndent(rh, "", "\t"); err != nil {
			log.Fatal("Can't save to .sainit file")
		}

		tools.Save2File(brh, path.Join(homedir, ".sainit"))

	} else {
		prh := &Roothome{}
		if err = json.Unmarshal(d, prh); err != nil {
			log.Fatal("Cant recover home dir")
		}

		savedir = prh.Rootdir

	}
	if savedir == "" {
		return nil
	}
	return &SARootConfig{HomeDir: savedir, CfgDir: path.Join(savedir, "config"), CfgFileName: "sa.json"}
}

func DefaultInitRootConfig(hdir string, force bool) *SARootConfig {
	var sar *SARootConfig
	if force {
		sar = forceInitRootConfig(hdir)
	} else {
		sar = unforceInitRootConfig(hdir)
	}

	//log.Println("Config Root:", sar.HomeDir)

	return sar
}

func DefaultInitConfig() *SAConfig {
	sa := &SAConfig{}

	sa.BootstrapIPAddress = []string{"103.45.98.72:50811", "174.7.124.45:50811"}
	sa.ConsumePath = "/public/keys/consume"
	sa.DownloadPath = "/download"
	sa.EmailPath = "/public/key/refresh"

	sa.UpdateClientSoftwarePath = "/public/app"
	sa.Uploadpath = "/upload"
	sa.TestIPAddressPath = "/localipaddress"
	sa.ListIpsPath = "/public/servers/list"
	sa.PostSocks5Path = "/postsocks5"
	sa.KeyDir = "key"
	sa.PidDir = "piddir"
	sa.FileDBDir = "filedb"
	sa.HomeIPDBFile = "homeipdbaddress"
	sa.FileStoreDir = "filestore"
	sa.VerifyPath = "/public/keys/verify"
	sa.RemoteServerIP = "207.148.9.49"
	sa.RemoteServerPort = 80
	sa.HttpListenPort = 50810
	sa.CmdListenIP = "127.0.0.1"
	sa.CmdListenPort = 50811
	sa.DhtListenPort = 50811
	sa.ListenTyp = "tcp4"
	sa.StaticFileDir = "staticfile"
	sa.LoginDir = "login"
	sa.Loginfile = "login.gptl"
	sa.LicenseAdminUser = [][]string{{"sofaadmin","J1jdNR8vQb"},{"nbsadmin","Dkf44u3Ad8"},}

	return sa
}

func (sar *SARootConfig) LoadCfg() *SAConfig {

	data, err := tools.OpenAndReadAll(path.Join(sar.CfgDir, sar.CfgFileName))

	if err != nil {
		return nil
	}

	cfg := &SAConfig{}

	if err = json.Unmarshal(data, cfg); err != nil {
		return nil
	}

	sar.SacInst = cfg

	cfg.Root = sar

	return cfg

}

func (sar *SARootConfig) IsInitialized() bool {
	cfgname := path.Join(sar.CfgDir, sar.CfgFileName)
	if !tools.FileExists(cfgname) {
		return false
	}

	return true
}

func (sar *SARootConfig) InitConfig(force bool) *SARootConfig {
	if sar.HomeDir == "" || sar.CfgDir == "" || sar.CfgFileName == "" {
		log.Fatal("Please Set Config Path")
	}

	if force {
		os.RemoveAll(sar.HomeDir)
	}

	cfgname := path.Join(sar.CfgDir, sar.CfgFileName)

	if !tools.FileExists(cfgname) {
		if !tools.FileExists(sar.CfgDir) {
			os.MkdirAll(sar.CfgDir, 0755)
		}
		sac := DefaultInitConfig()
		sar.SacInst = sac
		bjson, err := json.MarshalIndent(*sac, "", "\t")
		if err != nil {
			log.Fatal("Json module error")
		}
		tools.Save2File(bjson, cfgname)
	} else {
		bjson, err := tools.OpenAndReadAll(cfgname)
		if err != nil {
			log.Fatal(err)
		}

		sac := &SAConfig{}
		err = json.Unmarshal(bjson, sac)
		if err != nil {
			log.Fatal(err)
		}
		sar.SacInst = sac
	}

	filedbdir :=""

	if sar.SacInst.FileDBDir[0] == '/'{
		filedbdir = sar.SacInst.FileDBDir
	}else{
		filedbdir = path.Join(sar.HomeDir, sar.SacInst.FileDBDir)
	}

	filestoredir := ""
	if sar.SacInst.FileStoreDir[0] == '/'{
		filestoredir = sar.SacInst.FileStoreDir
	}else{
		filestoredir = path.Join(sar.HomeDir, sar.SacInst.FileStoreDir)
	}

	keydir := path.Join(sar.HomeDir, sar.SacInst.KeyDir)

	piddir := ""
	if sar.SacInst.PidDir[0] == '/'{
		piddir = sar.SacInst.PidDir
	}else{
		piddir = path.Join(sar.HomeDir, sar.SacInst.PidDir)
	}

	staticfile:=path.Join(sar.HomeDir,sar.SacInst.StaticFileDir)

	if !tools.FileExists(filedbdir) {
		os.MkdirAll(filedbdir, 0755)
	}
	if !tools.FileExists(filestoredir) {
		os.MkdirAll(filestoredir, 0755)
	}
	if !tools.FileExists(keydir) {
		os.MkdirAll(keydir, 0755)
	}
	if !tools.FileExists(piddir) {
		os.MkdirAll(piddir, 0755)
	}

	if !tools.FileExists(staticfile) {
		os.MkdirAll(staticfile, 0755)
	}

	loginfilename:=path.Join(staticfile,sar.SacInst.Loginfile)

	if !tools.FileExists(loginfilename){
		htmlfile.NewLoginFile(loginfilename)
	}

	sar.SacInst.Root = sar

	return sar
}

func (sar *SARootConfig) InitRSAKey(force bool) *SARootConfig {

	rsakeypath := path.Join(sar.HomeDir, "key")

	if !nbscrypt.RsaKeyIsExists(rsakeypath) || force {
		priv, _ := nbscrypt.GenerateKeyPair(2048)
		if err := nbscrypt.Save2FileRSAKey(rsakeypath, priv); err != nil {
			log.Fatal(err)
		}
	}

	return sar
}

func (sac *SAConfig)GetPidDir() string  {
	if sac.PidDir[0] == '/'{
		return sac.PidDir
	}

	return  path.Join(sac.Root.HomeDir,sac.PidDir)
}

func (sac *SAConfig)GetFileDbDir() string  {
	if sac.FileDBDir[0] == '/'{
		return sac.FileDBDir
	}

	return  path.Join(sac.Root.HomeDir,sac.FileDBDir)
}

func (sac *SAConfig)GetFileStoreDir() string  {
	if sac.FileStoreDir[0] == '/'{
		return sac.FileStoreDir
	}

	return  path.Join(sac.Root.HomeDir,sac.FileStoreDir)
}

func (sac *SAConfig) GenNbsRsaAddr() {
	if sac.PrivKey == nil {
		log.Fatal(errors.New("No Private Key Found"))
	}

	pubkeybytes := x509.MarshalPKCS1PublicKey(&sac.PrivKey.PublicKey)

	s := sha256.New()
	s.Write(pubkeybytes)

	sum := s.Sum(nil)

	sac.NbsRsaAddr = "91" + base58.Encode(sum)
}

func (sar *SARootConfig) LoadRsaKey() {

	if sar.SacInst == nil {
		log.Fatal(errors.New("No config instance"))
	}

	rsakeypath := path.Join(sar.HomeDir, "key")

	priv, _, err := nbscrypt.LoadRSAKey(rsakeypath)
	if err != nil {
		log.Fatal(err)
	}

	sar.SacInst.PrivKey = priv

	if sar.SacInst.NbsRsaAddr == "" {
		sar.SacInst.GenNbsRsaAddr()
		bjson, err := json.MarshalIndent(*sar.SacInst, "", "\t")
		if err != nil {
			log.Fatal("Json module error")
		} else {
			tools.Save2File(bjson, path.Join(sar.CfgDir, sar.CfgFileName))
		}
	}
}

func CheckUserPassword(username string,password string) bool {
	sac:=GetSAConfig()

	for _,authpair:=range sac.LicenseAdminUser{
		if len(authpair) == 2{
			if username == authpair[0] && password == authpair[1]{
				return true
			}
		}
	}

	return false
}


