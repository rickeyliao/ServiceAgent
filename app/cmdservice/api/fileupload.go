package api

import (
	"context"
	"github.com/kprc/nbsnetwork/tools"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"net"
	"os"
	"sync"

	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type CmdFileUpLoad struct {
}

func (cfu *CmdFileUpLoad) Uploadfile(ctx context.Context, req *pb.Fileuploadreq) (*pb.DefaultResp, error) {
	if req.Filepath == "" || req.Hostip == "" {
		return encResp("Param error"), nil
	}

	hiparr := strings.Split(req.Hostip, ":")
	if len(hiparr) != 2 {

		return encResp("host ip address error"), nil
	}

	if net.ParseIP(hiparr[0]) == nil {

		return encResp("host ip address error"), nil
	}

	if rport, err := strconv.Atoi(hiparr[1]); err != nil {

		return encResp("host ip address error"), nil
	} else {
		if rport < 1024 || rport > 65535 {

			return encResp("host ip address error"), nil
		}
	}

	if !tools.FileExists(req.Filepath) {
		return encResp("file not found"), nil
	}

	if hv, err := uploadfile(req.Hostip, req.Filepath); err != nil {
		return encResp("upload failed"), nil
	} else {
		msg := fmt.Sprintf("Upload file: %s\r\nTo host: %s\r\nRename To HashCode: %s", req.Filepath, req.Hostip, hv)
		return encResp(msg), nil
	}
}

func uploadfile(hostip string, filepath string) (string, error) {

	resultchan := make(chan error, 1)

	wg := &sync.WaitGroup{}
	piper, pipew := io.Pipe()
	defer piper.Close()

	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	hv := genrsa256hash(file)

	file.Close()

	m := multipart.NewWriter(pipew)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer pipew.Close()

		part, err := m.CreateFormFile("FileHash", hv)
		if err != nil {
			resultchan <- err
			return
		}
		defer m.Close()

		file, err := os.Open(filepath)
		if err != nil {
			resultchan <- err
			return
		}
		defer file.Close()

		if _, err = io.Copy(part, file); err != nil {
			resultchan <- err
			return
		}
		resultchan <- err
		return
	}()

	posturl := "http://" + hostip + "/upload"

	client := http.Client{Transport: &(http.Transport{DisableKeepAlives: true})}

	request, _ := http.NewRequest("POST", posturl, piper)
	request.Header.Set("Content-Type", m.FormDataContentType())

	resp, err := client.Do(request)

	//resp,err:=http.Post(posturl, m.FormDataContentType(), piper)
	wg.Wait()

	select {
	case err = <-resultchan:
		//todo...
	default:
	}

	if err != nil {
		fmt.Println(err)
	}

	if resp != nil && resp.Body != nil {
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}

	if err != nil {
		return "", err
	} else {
		return hv, nil
	}
}

func genrsa256hash(f *os.File) string {
	s := sha256.New()

	if _, err := io.Copy(s, f); err != nil {
		return ""
	}

	hv := s.Sum(nil)

	return "c1" + base58.Encode(hv)
}
