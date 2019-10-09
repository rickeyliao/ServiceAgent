package service

import (
	"bytes"
	"context"
	"crypto/rsa"
	"github.com/rickeyliao/ServiceAgent/agent/email"
	"github.com/rickeyliao/ServiceAgent/agent/key"
	"github.com/rickeyliao/ServiceAgent/agent/listallip"
	"github.com/rickeyliao/ServiceAgent/agent/software"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/ServiceAgent/service/checkip"
	"github.com/rickeyliao/ServiceAgent/service/file"
	"github.com/rickeyliao/ServiceAgent/service/license"
	"github.com/rickeyliao/ServiceAgent/service/localaddress"
	"github.com/rickeyliao/ServiceAgent/service/login"
	"github.com/rickeyliao/ServiceAgent/service/postsocks5"
	"github.com/rickeyliao/ServiceAgent/service/pubkey"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/tools/crypt/nbscrypt"
	"io/ioutil"
)

var (
	httpserver *http.Server
	quit       chan int
	wg         sync.WaitGroup
)

func Run(cfg *common.SAConfig) {

	remotehost := cfg.RemoteServerIP
	remoteport := cfg.RemoteServerPort

	if remotehost == "" && remoteport == 0 {
		log.Println("Please set remote host and port")
		return
	}

	mux := http.NewServeMux()

	ips := remotehost + ":" + strconv.Itoa(int(remoteport))
	common.NewRemoteUrl1(ips)

	mux.Handle(cfg.VerifyPath, key.NewKeyAuth())
	mux.Handle(cfg.ConsumePath, key.NewKeyImport())
	mux.Handle(cfg.EmailPath, email.NewEmailRecord())
	mux.Handle(cfg.UpdateClientSoftwarePath, software.NewUpdateSoft())
	mux.Handle(cfg.TestIPAddressPath, localaddress.NewLocalAddress())
	mux.Handle(cfg.ListIpsPath, listallip.NewListAllIps())
	mux.Handle(cfg.PostSocks5Path, postsocks5.NewPostSocks5())
	mux.Handle(cfg.Uploadpath, file.NewFileUpLoad())
	mux.Handle(cfg.DownloadPath, file.NewFileDownLoad())
	mux.Handle(cfg.LoginPath, login.NewLoginInfo())
	mux.Handle(cfg.CheckIPPath, checkip.NewCheckPrivateIP())
	mux.Handle(cfg.PubkeyPath, pubkey.NewHttpPubKey())

	listenportstr := ":" + strconv.Itoa(int(cfg.HttpListenPort))

	log.Println("Remote Server:", common.GetRemoteUrlInst().GetHostName(""))
	log.Println("Server Listen at:", listenportstr)
	log.Println("LocalNbsAddress:", cfg.NbsRsaAddr)

	if !cfg.IsCoordinator {
		quit = make(chan int, 0)
		wg.Add(1)
		go reportAddress()
	}

	if cfg.IsCoordinator {
		go localaddress.IntervalSave()
	}

	go license.IntervalSave()

	httpserver = &http.Server{Addr: listenportstr, Handler: mux}

	log.Fatal(httpserver.ListenAndServe())
}

func Stop() {
	if !common.GetSAConfig().IsCoordinator {
		quit <- 1
		wg.Wait()
	}

	if common.GetSAConfig().IsCoordinator {
		localaddress.Destroy()
	}

	license.Destroy()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	httpserver.Shutdown(ctx)
}

func getCryptSSInfo(pk *rsa.PublicKey) string {
	ssr := localaddress.GetSSReport()

	bssr, err := json.Marshal(*ssr)
	if err != nil {
		log.Println(err)
		return ""
	}

	encdata, err := nbscrypt.EncryptRSA(bssr, pk)
	if err != nil {
		log.Println(err, pk.Size()-11)
		return ""
	}

	//for debug
	//fmt.Println("for send crypt ssinfo :", base58.Encode(encdata))

	return base58.Encode(encdata)

}

func report(address string, ra *rsaaddr) {
	tp := http.Transport{DisableKeepAlives: true}
	c := &http.Client{Transport: &tp}

	if req, err := http.NewRequest("GET", "http://"+address+common.GetSAConfig().TestIPAddressPath, nil); err != nil {
		log.Println(err)
		return
	} else {

		req.Header.Add("nbsaddress", common.GetSAConfig().NbsRsaAddr)
		ips := common.GetAllLocalIpAddr()
		req.Header.Add("nataddrs", localaddress.LocalIPArr2string(ips))
		req.Header.Add("hostname", localaddress.GetMachineName())
		//req.Header.Add("nationality", strconv.Itoa(int(common.GetSAConfig().Nationality)))
		req.Header.Add("ssrinfo", getCryptSSInfo(ra.pk))

		if resp, errresp := c.Do(req); errresp != nil {
			log.Println(errresp)
			return
		} else {
			resp.Body.Close()
			//log.Println(resp)
		}

	}
}

type rsaaddr struct {
	pk      *rsa.PublicKey
	addr    string
	nbsaddr string
	ts      int64
	failcnt int32
}

func (ra *rsaaddr) print() {
	bpk := x509.MarshalPKCS1PublicKey(ra.pk)
	base58pk := base58.Encode(bpk)

	fmt.Println("rsaaddr.pk", base58pk)
	fmt.Println("rsaaddr.addr", ra.addr)
	fmt.Println("rsaaddr.nbsaddr", ra.nbsaddr)
	fmt.Println("rsaaddr.ts", strconv.FormatInt(ra.ts, 10))
	fmt.Println("rsaaddr.failcnt", strconv.Itoa(int(ra.failcnt)))
}

func reqrsaaddr(addr string) *rsaaddr {
	tp := http.Transport{DisableKeepAlives: true}
	c := &http.Client{Transport: &tp}

	var ra *rsaaddr

	r := bytes.NewReader([]byte(pubkey.GetNbsPubkey()))

	if req, err := http.NewRequest("POST", "http://"+addr+common.GetSAConfig().PubkeyPath, r); err != nil {
		log.Println(err)
		return nil
	} else {
		if resp, errresp := c.Do(req); errresp != nil {
			log.Println(errresp)
			return nil
		} else {

			if pkjson, err := ioutil.ReadAll(resp.Body); err == nil {
				nbsaddr, rsapk := pubkey.UnMarshalPubKey(pkjson)
				ra = &rsaaddr{}
				ra.pk = rsapk
				ra.addr = addr
				ra.nbsaddr = nbsaddr
				ra.ts = tools.GetNowMsTime() / 1000
			} else {
				log.Println(err)
			}

			resp.Body.Close()
		}
	}

	return ra
}

func updatemapaddr(addr string, mapaddr map[string]*rsaaddr) *rsaaddr {

	if addr == "" {
		return nil
	}
	var v *rsaaddr
	var ok bool

	now := tools.GetNowMsTime() / 1000

	if v, ok = mapaddr[addr]; ok {
		if v.failcnt == 0 && now-v.ts < 86400 {
			return v
		}
	}
	ra := reqrsaaddr(addr)

	if ra == nil {
		if v != nil {
			v.failcnt++
		}
		return nil
	} else {
		mapaddr[addr] = ra
	}

	return ra
}

func reportAddress() {
	var count int64

	mapaddr := make(map[string]*rsaaddr)

	for {
		count++
		if count%300 == 0 {
			for _, addr := range common.GetSAConfig().ReportServerIPAddress {
				ra := updatemapaddr(addr, mapaddr)
				if ra != nil {
					report(addr, ra)
				}
				time.Sleep(time.Second * 1)
			}
		}
		time.Sleep(time.Second * 1)
		select {
		case <-quit:
			wg.Done()
			return
		default:
			//todo...
		}
	}
}
