#!/usr/bin/env bash

function upload(){
nameathost=$1
exec=$2

sftp $nameathost <<EOF
put bin/$exec.gz
put config.json
EOF
}




function deploy(){
nameathost=$1
exec=$2

ssh $nameathost <<EOF
ps -ef|grep $exec|grep -v grep|awk '{print \$2}'
ps -ef|grep $exec|grep -v grep|awk '{print \$2}'|sudo xargs kill
rm -f ./$exec
gzip -d $exec.gz
sudo nohup ./$exec >$exec.out 2>&1 &
date
sleep 2
cat ./$exec.out
ps -ef|grep $exec|grep -v grep|awk '{print \$2}'
EOF

}

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

ROOT=`pwd`
echo $ROOT

#upload root@103.19.2.78  ddns-server-linux64
#deploy root@103.19.2.78  ddns-server-linux64


#upload pi@192.168.31.83  ddns-client-linuxarm7
#deploy pi@192.168.31.83  ddns-client-linuxarm7

upload pi@192.168.50.168  ddns-client-linuxarm7
deploy pi@192.168.50.168  ddns-client-linuxarm7
echo "done"