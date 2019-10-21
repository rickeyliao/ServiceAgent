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
	"github.com/rickeyliao/ServiceAgent/htmlfile"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

type SAConfig struct {
	DownloadPath             string     `json:"downloadpath"`
	Uploadpath               string     `json:"uploadpath"`
	TestIPAddressPath        string     `json:"testipaddress"`
	PostSocks5Path           string     `json:"postsocks5path"`
	VerifyPath               string     `json:"verifypath"`
	ConsumePath              string     `json:"consumepath"`
	ListIpsPath              string     `json:"listipspath"`
	ServerDelete             string     `json:"serverdel"`
	ServerAdd                string     `json:"serveradd"`
	EmailPath                string     `json:"emailpath"`
	UpdateClientSoftwarePath string     `json:"updateclientsoftwarepath"`
	PubkeyPath               string     `json:"pubkeypath"`
	KeyDir                   string     `json:"keydir"`
	PidDir                   string     `json:"piddir"`
	FileDBDir                string     `json:"filedbdir"`
	HomeIPDBFile             string     `json:"homeipdbfile"`
	LicenseDBFile            string     `json:"licensedbfile"`
	FileStoreDB              string     `json:"filestoredb"`
	FileStoreDir             string     `json:"filestoredir"`
	RemoteServerIP           string     `json:"remoteserverip"`
	RemoteServerPort         uint16     `json:"remoteserverport"`
	HttpListenPort           uint16     `json:"httplistenport"` //50810 tcp/http for file transfer
	BootstrapIPAddress       []string   `json:"bootstrapipaddress"`
	ReportServerIPAddress    []string   `json:"reportserveraddress"`
	ListenTyp                string     `json:"listentyp"`
	NbsRsaAddr               string     `json:"nbsaddr"`
	CmdListenIP              string     `json:"cmdlistenip"`
	CmdListenPort            uint16     `json:"cmdlistenport"` //50811 tcp for cmd
	DhtListenPort            uint16     `json:"dhtlistenport"` //50811 udp for control message
	WebServerPort            uint16     `json:"webserverport"` //50814
	StaticFileDir            string     `json:"staticfiledir"`
	LoginPath                string     `json:"logindir"`
	Loginfile                string     `json:"loginfile"`
	CheckIPPath              string     `json:"checkipdir"`
	CheckIPFile              string     `json:"checkipfile"`
	LicenseAdminUser         [][]string `json:"licenseadminuser"`
	ShadowSockServerSwitch   bool       `json:"shadowsockserverswitch"`
	ShadowSockPort           uint16     `json:"shadowsockport"` //50812
	ShadowSockPasswd         string     `json:"sspasswd"`
	ShadowSockMethod         string     `json:"ssmethod"`
	ShadowSockStatFile       string     `json:"ssstatefile"`
	StatisticDir             string     `json:"ssdir"`
	HostName                 string     `json:"hostname"`
	IsCoordinator            bool       `json:"iscoordinator"`
	Nationality              int32      `json:"nationality"`
	Version                  string		`json:"-"`

	CoinBase                 bool		`json:"coinbase"`
	CoinCount                float64    `json:"coinCount"`
	TrafficCnt               int64 		`json:"trafficcnt"`


	PrivKey *rsa.PrivateKey `json:"-"`
	Root    *SARootConfig   `json:"-"`
}

type SARootConfig struct {
	HomeDir     string
	CfgDir      string
	CfgFileName string
	SacInst     *SAConfig
	needSave    bool
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

	cfgdir := path.Join(sahome, "config")

	return &SARootConfig{HomeDir: sahome, CfgDir: cfgdir, CfgFileName: "sa.json", needSave: true}
}

func (sac *SAConfig)SetNBSVersion(ver string)  {

	sac.Version = ver
}

func unforceInitRootConfig(hdir string) *SARootConfig {
	var homedir string
	var savedir string
	var d []byte
	var err error
	var nds bool

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

		nds = true

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
	return &SARootConfig{HomeDir: savedir, CfgDir: path.Join(savedir, "config"), CfgFileName: "sa.json", needSave: nds}
}

func DefaultInitRootConfig(hdir string, force bool) *SARootConfig {
	var sar *SARootConfig
	if force {
		sar = forceInitRootConfig(hdir)
	} else {
		sar = unforceInitRootConfig(hdir)
	}

	if sar != nil && sar.needSave {
		sar.Save()
	}

	return sar
}

func DefaultInitConfig() *SAConfig {
	sa := &SAConfig{}

	sa.BootstrapIPAddress = []string{"103.45.98.72:50811", "24.86.164.242:50811"}
	sa.ReportServerIPAddress = []string{"103.45.98.72:50810", "24.86.164.242:50810"}
	sa.ConsumePath = "/public/keys/consume"
	sa.DownloadPath = "/download"
	sa.EmailPath = "/public/key/refresh"
	sa.UpdateClientSoftwarePath = "/public/app"
	sa.Uploadpath = "/upload"
	sa.TestIPAddressPath = "/localipaddress"
	sa.ListIpsPath = "/public/servers/list"
	sa.ServerAdd = "/public/servers/add"
	sa.ServerDelete = "/public/servers/delete"
	sa.PostSocks5Path = "/postsocks5"
	sa.PubkeyPath = "/pubkey/fetch"
	sa.KeyDir = "key"
	sa.PidDir = "piddir"
	sa.FileDBDir = "filedb"
	sa.HomeIPDBFile = "homeipdbaddress"
	sa.FileStoreDB = "hashfilestoredb"
	sa.LicenseDBFile = "licensedb"
	sa.FileStoreDir = "filestore"

	sa.VerifyPath = "/public/keys/verify"
	sa.RemoteServerIP = "207.148.9.49"
	sa.RemoteServerPort = 80
	sa.HttpListenPort = 50810
	sa.CmdListenIP = "127.0.0.1"
	sa.CmdListenPort = 50811
	sa.DhtListenPort = 50811
	sa.WebServerPort = 50814
	sa.ListenTyp = "tcp4"
	sa.StaticFileDir = "staticfile"
	sa.LoginPath = "/login"
	sa.CheckIPPath = "/checkip"
	sa.Loginfile = "login.gptl"
	sa.CheckIPFile = "checkip.gptl"
	sa.ShadowSockServerSwitch = false
	sa.ShadowSockPort = 50812
	sa.ShadowSockStatFile = "shadowsock.stat"
	sa.StatisticDir = "stat"

	sa.ShadowSockPasswd = ""
	sa.ShadowSockMethod = ""
	sa.HostName = ""
	sa.CoinBase = true
	sa.IsCoordinator = false
	sa.LicenseAdminUser = [][]string{{"sofaadmin", "J1jdNR8vQb"}, {"nbsadmin", "Dkf44u3Ad8"}}
	sa.Nationality = 1 //1 American 86 China

	return sa
}

func (sar *SARootConfig) Save() *SARootConfig {

	if sar.HomeDir == "" {
		return sar
	}

	//save homedir to .sainit file
	rh := Roothome{sar.HomeDir}

	if brh, err := json.MarshalIndent(rh, "", "\t"); err != nil {
		log.Fatal("Can't save to .sainit file")
	} else {
		if homedir, err1 := tools.Home(); err1 != nil {
			log.Fatal("Can't Get Home Directory")
		} else {
			tools.Save2File(brh, path.Join(homedir, ".sainit"))
		}
	}

	return sar
}

func (sar *SARootConfig) LoadCfg() *SAConfig {

	data, err := tools.OpenAndReadAll(path.Join(sar.CfgDir, sar.CfgFileName))

	if err != nil {
		return nil
	}

	cfg := DefaultInitConfig()

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

const (
	InitNone  int = 0
	InitTrue  int = 1
	InitFalse int = 2
)

type ConfigInitParam struct {
	Force       bool
	IsCoord     int
	Hostname    string
	SS          string
	Nationality int32
}

func (sar *SARootConfig) InitConfig(cip *ConfigInitParam) *SARootConfig {
	var nds bool
	if sar.HomeDir == "" || sar.CfgDir == "" || sar.CfgFileName == "" {
		log.Fatal("Please Set Config Path")
	}

	if cip.Force {
		os.RemoveAll(sar.HomeDir)
		nds = true
	}

	cfgname := path.Join(sar.CfgDir, sar.CfgFileName)

	if !tools.FileExists(cfgname) {
		if !tools.FileExists(sar.CfgDir) {
			os.MkdirAll(sar.CfgDir, 0755)
		}
		sac := DefaultInitConfig()
		sar.SacInst = sac
		nds = true
	} else {
		bjson, err := tools.OpenAndReadAll(cfgname)
		if err != nil {
			log.Fatal(err)
		}

		sac := DefaultInitConfig()
		err = json.Unmarshal(bjson, sac)
		if err != nil {
			log.Fatal(err)
		}
		sar.SacInst = sac
	}

	if cip.Hostname != "" {
		if sar.SacInst.HostName != cip.Hostname {
			sar.SacInst.HostName = cip.Hostname
			nds = true
		}
	}

	if cip.IsCoord != InitNone {
		iscoordb := false
		if cip.IsCoord == InitTrue {
			iscoordb = true
		}
		if sar.SacInst.IsCoordinator != iscoordb {
			sar.SacInst.IsCoordinator = iscoordb
			nds = true
		}
	}

	if cip.Nationality > 0 {
		sar.SacInst.Nationality = cip.Nationality
		nds = true
	}

	//if cip.Location != sar.SacInst.Location {
	//	sar.SacInst.Location = cip.Location
	//	nds = true
	//}

	filedbdir := ""

	if sar.SacInst.FileDBDir[0] == '/' {
		filedbdir = sar.SacInst.FileDBDir
	} else {
		filedbdir = path.Join(sar.HomeDir, sar.SacInst.FileDBDir)
	}

	filestoredir := ""
	if sar.SacInst.FileStoreDir[0] == '/' {
		filestoredir = sar.SacInst.FileStoreDir
	} else {
		filestoredir = path.Join(sar.HomeDir, sar.SacInst.FileStoreDir)
	}

	keydir := path.Join(sar.HomeDir, sar.SacInst.KeyDir)

	piddir := ""
	if sar.SacInst.PidDir[0] == '/' {
		piddir = sar.SacInst.PidDir
	} else {
		piddir = path.Join(sar.HomeDir, sar.SacInst.PidDir)
	}

	staticfile := path.Join(sar.HomeDir, sar.SacInst.StaticFileDir)

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

	loginfilename := path.Join(staticfile, sar.SacInst.Loginfile)

	if !tools.FileExists(loginfilename) {
		htmlfile.NewLoginFile(loginfilename)
	}

	checkipfilename := path.Join(staticfile, sar.SacInst.CheckIPFile)
	if !tools.FileExists(checkipfilename) {
		htmlfile.NewCheckIPFile(checkipfilename)
	}

	statdir := path.Join(sar.HomeDir, sar.SacInst.StatisticDir)
	if !tools.FileExists(statdir) {
		os.MkdirAll(statdir, 0755)
	}

	sar.SacInst.Root = sar

	if nds {
		sar.SacInst.Save()
	}

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

func (sac *SAConfig) GetPidDir() string {
	if sac.PidDir[0] == '/' {
		return sac.PidDir
	}

	return path.Join(sac.Root.HomeDir, sac.PidDir)
}

func (sac *SAConfig) GetFileDbDir() string {
	if sac.FileDBDir[0] == '/' {
		return sac.FileDBDir
	}

	return path.Join(sac.Root.HomeDir, sac.FileDBDir)
}

func (sac *SAConfig) GetSSStatFile() string {
	return path.Join(sac.Root.HomeDir, sac.StatisticDir, sac.ShadowSockStatFile)
}

func (sac *SAConfig) GetFileStoreDir() string {
	if sac.FileStoreDir[0] == '/' {
		return sac.FileStoreDir
	}

	return path.Join(sac.Root.HomeDir, sac.FileStoreDir)
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

func (sac *SAConfig) GetSSPasswd() string {

	if sac.ShadowSockPasswd == "" {
		return ""
	}

	if encpasswd, err := base58.Decode(sac.ShadowSockPasswd); err != nil {
		return ""
	} else {

		if data, errd := nbscrypt.DecryptRsa(encpasswd, sac.PrivKey); errd != nil {
			return ""
		} else {
			return string(data)
		}
	}

}

func (sac *SAConfig) GetSSMethod() string {
	if sac.ShadowSockMethod == "" {
		return ""
	}

	if encmethod, err := base58.Decode(sac.ShadowSockMethod); err != nil {
		return ""
	} else {

		if data, errd := nbscrypt.DecryptRsa(encmethod, sac.PrivKey); errd != nil {
			return ""
		} else {
			return string(data)
		}
	}
}

func (sac *SAConfig) GetPubKey() string {
	pubkeybytes := x509.MarshalPKCS1PublicKey(&sac.PrivKey.PublicKey)

	return base58.Encode(pubkeybytes)
}

func (sac *SAConfig) Save() {
	bjson, err := json.MarshalIndent(*sac, "", "\t")
	if err != nil {
		log.Fatal("Json module error")
	} else {
		tools.Save2File(bjson, path.Join(sac.Root.CfgDir, sac.Root.CfgFileName))
	}
}

func (sar *SARootConfig) SetShadowSockParam(param string) {
	parr := strings.Split(param, ":")

	if len(parr) != 2 {
		log.Println("Set shadowsock error, use default parameter")
		return
	}

	encpasswd, err := nbscrypt.EncryptRSA([]byte(parr[0]), &sar.SacInst.PrivKey.PublicKey)
	if err != nil {
		log.Println("Internal error")
	}
	passwd := base58.Encode(encpasswd)
	sar.SacInst.ShadowSockPasswd = passwd

	encmethod, err := nbscrypt.EncryptRSA([]byte(parr[1]), &sar.SacInst.PrivKey.PublicKey)
	if err != nil {
		log.Println("Internal error")
	}

	method := base58.Encode(encmethod)
	sar.SacInst.ShadowSockMethod = method

	sar.SacInst.ShadowSockServerSwitch = true

	sar.SacInst.Save()
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
		sar.SacInst.Save()
	}
}

func CheckUserPassword(username string, password string) bool {
	sac := GetSAConfig()

	for _, authpair := range sac.LicenseAdminUser {
		if len(authpair) == 2 {
			if username == authpair[0] && password == authpair[1] {
				return true
			}
		}
	}

	return false
}
