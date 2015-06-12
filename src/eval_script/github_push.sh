#!/bin/bash

i=$1
no=$2
echo "no is"
echo $no

while [ $i -lt $no ]
do
    dd if=/dev/zero of=file_to-create-$i bs=2k count=$3
    git add file_to-create-$i

    git commit -m "done"

    git pull

    time git push
    
    i=$[$i+1]
done


