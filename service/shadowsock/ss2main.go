package shadowsock

import (
	"github.com/pkg/errors"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/shadowsocks/go-shadowsocks2/core"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var config struct {
	Verbose    bool
	UDPTimeout time.Duration
}

var (
	SSTCPListener   *net.Listener
	SSUDPPacketConn *net.PacketConn
)

//var logger = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

func logf(f string, v ...interface{}) {
	if config.Verbose {
		log.Println(v)
	}
}

func ss2server(port int, passwd, method string) error {

	if passwd == "" {
		return errors.New("Please Set Passwd")
	}

	var flags struct {
		Server   string
		Cipher   string
		Password string
	}

	flags.Cipher = "AES-256-CFB"
	if method != "" {
		flags.Cipher = strings.ToUpper(method)
	}

	flags.Password = passwd

	if port > 1024 {
		flags.Server = ":" + strconv.Itoa(port)
	} else {
		flags.Server = ":50812"
	}

	var key []byte

	config.Verbose = false

	addr := flags.Server
	cipher := flags.Cipher
	password := flags.Password

	ciph, err := core.PickCipher(cipher, key, password)
	if err != nil {
		log.Fatal(err)
	}

	go udpRemote(addr, ciph.PacketConn)

	tcpRemote(addr, ciph.StreamConn)

	return nil
}

func parseURL(s string) (addr, cipher, password string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}

	addr = u.Host
	if u.User != nil {
		cipher = u.User.Username()
		password, _ = u.User.Password()
	}
	return

}

func StartSS2Server() {
	sa := common.GetSAConfig()
	ss2server(int(sa.ShadowSockPort), sa.GetSSPasswd(), sa.GetSSMethod())
}

func StopSS2Server() {
	if SSTCPListener != nil {
		lis := *SSTCPListener
		lis.Close()
	}

	if SSUDPPacketConn != nil {
		sup := *SSUDPPacketConn
		sup.Close()
	}
}
