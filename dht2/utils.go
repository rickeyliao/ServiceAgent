package dht2

import (
	"net"
	"strings"
	"math/big"
	"encoding/binary"
	"github.com/pkg/errors"
)

func GetAllLocalIps() []net.IP {
	ips := make([]net.IP, 0)
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil && !strings.Contains(ipnet.IP.String(), "169.254") &&
					!strings.Contains(ipnet.IP.String(), "127.0.0.1") {
					ips = append(ips, ipnet.IP)
				}
			}
		}

	}

	return ips
}

func NbsXorUInt(x,y uint32) (*big.Int,error)  {
	bx,by:=make([]byte,4),make([]byte,4)

	binary.BigEndian.PutUint32(bx,x)
	binary.BigEndian.PutUint32(by,y)

	return NbsXor(bx,by)
}


func NbsXor(x []byte,y []byte) (*big.Int,error) {

	bgx:=(&big.Int{}).SetBytes(x)
	bgy:=(&big.Int{}).SetBytes(y)


	z:=&big.Int{}

	z.Xor(bgx,bgy)

	return z,nil
}


var xorlen =[...]int{
	0, 1, 2, 2,3, 3, 3, 3, 4,4,4, 4, 4, 4, 4,4, //15
	5, 5, 5, 5,5, 5, 5, 5, 5,5, 5, 5, 5, 5,5, 5, //31
	6, 6, 6,6, 6, 6, 6, 6,6, 6, 6, 6, 6,6, 6, 6, //47
	6, 6,6, 6, 6, 6, 6,6, 6, 6, 6, 6,6, 6, 6, 6, //63
	7,7, 7, 7, 7, 7,7, 7, 7, 7, 7,7, 7, 7, 7, 7,
	7, 7, 7, 7, 7,7, 7, 7, 7, 7,7, 7, 7, 7, 7,7,
	7, 7, 7, 7,7, 7, 7, 7, 7,7, 7, 7, 7, 7,7, 7,
	7, 7, 7,7, 7, 7, 7, 7,7, 7, 7, 7, 7,7, 7, 7, //127
	8, 8,8, 8, 8, 8, 8,8, 8, 8, 8, 8,8, 8, 8, 8,
	8,8, 8, 8, 8, 8,8, 8, 8, 8, 8,8, 8, 8, 8, 8,
	8, 8, 8, 8, 8,8, 8, 8, 8, 8,8, 8, 8, 8, 8,8,
	8, 8, 8, 8,8, 8, 8, 8, 8,8, 8, 8, 8, 8,8, 8,
	8, 8, 8,8, 8, 8, 8, 8,8, 8, 8, 8, 8,8, 8, 8,
	8, 8,8, 8, 8, 8, 8,8, 8, 8, 8, 8,8, 8, 8, 8,
	8,8, 8, 8, 8, 8,8, 8, 8, 8, 8,8, 8, 8, 8, 8,
	8, 8, 8, 8, 8,8, 8, 8, 8, 8,8, 8, 8, 8, 8,8, //255
}

func NbsXorLen(x []byte,y[]byte) (int,error)  {

	if len(x) != len(y){
		return 0,errors.New("length not correct")
	}

	z:=make([]byte,len(x))

	for ix,bx:=range x{
		z[ix]=bx^y[ix]
	}

	var iz int
	var bz byte

	for iz,bz=range z{
		if bz != 0{
			break
		}
	}

	left := len(z) - iz -1

	return left*8 + xorlen[int(bz)],nil
}

func NbsXorUintLen(x,y uint32) (int, error)  {
	bx,by:=make([]byte,4),make([]byte,4)

	binary.BigEndian.PutUint32(bx,x)
	binary.BigEndian.PutUint32(by,y)

	return NbsXorLen(bx,by)
}

func NbsBigIntLen(x *big.Int) int  {
	z:=x.Bytes()


	var iz int
	var bz byte

	for iz,bz=range z{
		if bz != 0{
			break
		}
	}

	left := len(z) - iz -1

	return left*8 + xorlen[int(bz)]
}