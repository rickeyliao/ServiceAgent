#/bin/sh

go get -u -v github.com/inconshreveable/mousetrap

go get -u -v github.com/elazarl/go-bindata-assetfs

go get -u -v github.com/jteeuwen/go-bindata/...

make

make -f MakefileAMD64Lnx

make -f MakefileArm

make -f MakefileAMD64Win

