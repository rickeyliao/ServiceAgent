package db

import (
	"bufio"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"reflect"
)

type filedbkev struct {
	Key   string `json:"k"`
	Value string `json:"v"`
}

type filedb struct {
	filepath string
	f        *os.File
	mkey     map[string]string
}

func NewFileDb(filepath string) *filedb {
	return &filedb{filepath: filepath, mkey: make(map[string]string)}
}

func (fdb *filedb) Open() NbsDbInter {
	if fdb.filepath == "" {
		log.Fatal("No Fill ")
	}

	flag := os.O_RDWR | os.O_APPEND

	if !tools.FileExists(fdb.filepath) {
		flag |= os.O_CREATE
	}
	if f, err := os.OpenFile(fdb.filepath, flag, 0755); err != nil {
		log.Fatal("Can't open file")
	} else {
		fdb.f = f
	}

	fdb.load()

	fdb.f.Close()

	fdb.f = nil

	return fdb
}

func (fdb *filedb) load() {
	if fdb.f == nil {
		return
	}
	bf := bufio.NewReader(fdb.f)

	for {
		if line, _, err := bf.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}

			if err == bufio.ErrBufferFull {
				log.Fatal("Buffer full")
				break
			}

			if len(line) > 0 {
				//pending drop it
				log.Fatal("Reading pending")
				break
			}

		} else {
			if len(line) > 0 {
				fdb.tomap(line)
			}
		}
	}
}

func (fdb *filedb) tomap(line []byte) {

	k := &filedbkev{}

	if err := json.Unmarshal(line, k); err != nil {
		return
	} else {
		fdb.Insert(k.Key, k.Value)
	}
}

func (fdb *filedb) Insert(key string, value string) error {
	if _, ok := fdb.mkey[key]; !ok {
		fdb.mkey[key] = value
	} else {
		return errors.New("Duplicate key")
	}

	return nil
}

func (fdb *filedb) Delete(key string) {
	delete(fdb.mkey, key)
}

func (fdb *filedb) Find(key string) (string, error) {
	if v, ok := fdb.mkey[key]; ok {
		return v, nil
	}
	return "", errors.New("Not Found")
}

func (fdb *filedb) Update(key string, value string) {

	fdb.mkey[key] = value
}

func (fdb *filedb) write(data []byte) {
	if fdb.f == nil || fdb.filepath == "" {
		flag := os.O_WRONLY | os.O_TRUNC

		if !tools.FileExists(fdb.filepath) {
			flag |= os.O_CREATE
		}
		if f, err := os.OpenFile(fdb.filepath, flag, 0755); err != nil {
			log.Fatal("Can't open file")
			return
		} else {
			fdb.f = f
		}
	}

	fdb.f.Write(data)
}

func (fdb *filedb) Save() {

	listkey := reflect.ValueOf(fdb.mkey).MapKeys()
	for _, key := range listkey {
		k := key.Interface().(string)

		fk := &filedbkev{}

		fk.Key = k
		fk.Value = fdb.mkey[k]

		if bj, err := json.Marshal(fk); err != nil {
			log.Println("save error", k, fk.Value)
		} else {
			fdb.write(bj)
		}

	}

	if fdb.f != nil {
		fdb.f.Close()
	}

	fdb.f = nil

}
