#!/bin/bash

no=$1
echo "no is"
echo $no
i="0"

while [ $i -lt $no ]
do
    pushd . &> /dev/null
    cd "$i"
    dd if=/dev/zero of=file_to-create-$i bs=10k count=1000
    go run $ZINGPATH/cmd/zing/main.go add file_to-create-$i

    go run $ZINGPATH/cmd/zing/main.go commit -m "done"

    go run $ZINGPATH/cmd/zing/main.go pull

    go run $ZINGPATH/cmd/zing/main.go push
    
    i=$[$i+1]
    port=$[$port+1]
    popd &> /dev/null
done


