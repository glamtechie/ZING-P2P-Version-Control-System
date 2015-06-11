#!/bin/bash

port=29000
no=$1
echo "no is"
echo $no
i="0"

while [ $i -lt $no ]
do
    pushd . &> /dev/null
    mkdir "$i"
    cd "$i"
    go run $ZINGPATH/cmd/zing/main.go clone $2  $port
    echo $2
    go run $ZINGPATH/cmd/zing-server/main.go $port  &
    sleep 1
    i=$[$i+1]
    port=$[$port+1]
    popd &> /dev/null
done


