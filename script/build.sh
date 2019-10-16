#!/bin/bash

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

ROOT=`pwd`
echo $ROOT
bindir=$ROOT/bin
mkdir -p $bindir

build() {
    local name
    local GOOS
    local GOARCH

    if [[ $1 == "darwin" ]]; then
        # Enable CGO for OS X so change network location will not cause problem.
        export CGO_ENABLED=1
    else
        export CGO_ENABLED=0
    fi

    prog=ddns_$4
    pushd src/cmd/$prog
    name=ddns-$4-$3
    echo "building $name"
    GOOS=$1 GOARCH=$2 GOARM=7 go build -a || exit 1

    if [[ $4 == "server" ]]; then
        rice append --exec $prog
    fi

    if [[ $1 == "windows" ]]; then
        mv $prog.exe $ROOT/script/
        pushd $ROOT/script/
        cp $ROOT/config.json sample-config.json
        cp $ROOT/sample-config/client-multi-server.json multi-server.json
        zip $name.zip $prog.exe shadowsocks.exe sample-config.json multi-server.json
        rm -f $prog.exe sample-config.json multi-server.json
        mv $name.zip $bindir
        popd
    else
        mv $prog $name
        gzip -f $name
        mv $name.gz $bindir
    fi
    popd
}

build darwin amd64 mac64 client
build linux amd64 linux64 client
build linux arm linuxarm7 client


build darwin amd64 mac64 server
build linux amd64 linux64 server

