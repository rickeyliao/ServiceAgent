package common

type SARootConfig struct {
	CfgDir      string
	CfgFileName string
	CfgFileType string
}


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
}










