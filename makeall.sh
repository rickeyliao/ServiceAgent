#/bin/sh

go get -u -v github.com/go-ole/go-ole/...

go get -u -v github.com/StackExchange/wmi

go get -u -v github.com/shirou/gopsutil/...

go get -u -v github.com/inconshreveable/mousetrap

go get -u -v github.com/elazarl/go-bindata-assetfs

go get -u -v github.com/jteeuwen/go-bindata/...

make

make -f MakefileAMD64Mac

make -f MakefileAMD64Lnx

make -f MakefileArm

make -f MakefileAMD64Win

