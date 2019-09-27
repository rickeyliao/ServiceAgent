package common

import (
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/nbsnetwork/tools"
	"io"
	"net/http"
	"os"
	"path"
)

func CheckNbsCotentHash(hv string) bool {
	if len(hv) < 2 {
		return false
	}
	if hv[:2] != "c1" {
		return false
	}

	hashv := hv[2:]
	v := base58.Decode(hashv)
	if v == nil || len(v) != 32 {
		return false
	}

	return true
}

func CheckNbsNodeHash(hv string) bool {
	if len(hv) < 2 {
		return false
	}
	if hv[:2] != "91" {
		return false
	}

	hashv := hv[2:]
	v := base58.Decode(hashv)
	if v == nil || len(v) != 32 {
		return false
	}

	return true
}

func GetSaveFilePath(filename string) string {
	sac := GetSAConfig()

	arrpath := getArrPath(filename)

	filepath := sac.GetFileStoreDir()
	for i := 0; i < len(arrpath); i++ {
		filepath = path.Join(filepath, arrpath[i])
	}

	if !tools.FileExists(filepath) {
		os.MkdirAll(filepath, 0755)
	}

	absfilename := path.Join(filepath, filename)

	return absfilename
}

func getArrPath(filename string) []string {
	arrpath := make([]string, 0)

	s := []byte(filename)

	for i := len(filename); i > 0; i = i - 2 {
		s := s[:i]
		if len(s) >= 2 {
			arrpath = append(arrpath, string(s[i-2:i]))

		} else {
			break
		}
		if len(arrpath) >= 4 {
			break
		}
	}

	if len(arrpath) > 0 {
		arrret := make([]string, 0)
		for i := len(arrpath) - 1; i >= 0; i-- {
			arrret = append(arrret, arrpath[i])
		}
		return arrret
	}

	return nil
}

func DownloadFile(hostip string, savepath string, filehash string) error {
	tp := http.Transport{DisableKeepAlives: true}
	c := &http.Client{Transport: &tp}

	geturl := "http://" + hostip + "/download"

	if req, err := http.NewRequest("GET", geturl, nil); err != nil {

		return err
	} else {

		req.Header.Add("FileHash", filehash)

		if resp, errresp := c.Do(req); errresp != nil {

			return errresp
		} else {
			defer resp.Body.Close()
			message := resp.Header.Get("message")
			if message == "FileNotFound" {
				return errors.New("File Not Found")
			}
			if savepath == "" {
				savepath = GetSAConfig().GetFileStoreDir()
			}
			f, err := os.OpenFile(path.Join(savepath, filehash), os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {
				return err
			}
			defer f.Close()
			io.Copy(f, resp.Body)

		}

		return nil

	}

}
