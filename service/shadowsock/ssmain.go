package shadowsock

import (
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/myshadowsock/cmd/SSServer"
)

func StartSSServer() {
	sa := common.GetSAConfig()

	SSServer.SSDaemon(int(sa.ShadowSockPort), sa.GetSSPasswd(), sa.GetSSMethod())
}

func StopSSServer() {
	SSServer.SSStop()
}
