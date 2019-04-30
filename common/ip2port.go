package common

import "hash/fnv"

func GetPort(ips string) uint16  {
	s := fnv.New64()
	s.Write([]byte(ips))
	h:=s.Sum64()

	p:= h & (0x3FF)

	p+=50000

	return uint16(p)
}
