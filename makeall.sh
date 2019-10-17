#/bin/sh

go get -u github.com/jteeuwen/go-bindata/...

make

make -f MakefileAMD64Lnx

make -f MakefileArm

make -f MakefileAMD64Win

