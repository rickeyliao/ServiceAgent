package transservice

import (
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/ServiceAgent/translayer"
	"net"
	"strconv"
)

type BigFileTransService struct {
	closechan chan error
}

func (bfts *BigFileTransService) StartService() error {

	return nil
}

func (bfts *BigFileTransService) UdpServerStart() (err error) {

	udpport := common.GetSAConfig().DhtListenPort

	s := ":" + strconv.Itoa(int(udpport))

	var pc net.PacketConn

	if pc, err = net.ListenPacket("udp4", s); err != nil {
		return
	}
	defer pc.Close()

	go func() {
		for {
			var n int
			var addr net.Addr
			buf := make([]byte, translayer.UPD_RCV_BUF_LEN)
			if n, addr, err = pc.ReadFrom(buf); err != nil {
				bfts.closechan <- err
				return
			}
			pc.SetWriteDeadline()

		}
	}()

	select {
	case err = <-bfts.closechan:
		return
	}

}
