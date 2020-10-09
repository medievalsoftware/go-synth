# /bin/sh

export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++
export AR=x86_64-w64-mingw32-ar

export CGO_ENABLED=1
export CGO_LDFLAGS='-static -s'

export GOOS=windows
export GOARCH=amd64
export GOEXE=".exe"

go build -v -x  .