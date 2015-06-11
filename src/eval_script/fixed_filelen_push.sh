#!/bin/bash

i=$2
no=$3
echo "no is"
echo $no

pushd . &> /dev/null
cd "$1"
while [ $i -lt $no ]
do
    dd if=/dev/zero of=file_to-create-$i bs=10k count=1000
    zing add file_to-create-$i

    zing commit -m "done"

    zing pull

    time zing push
    
    i=$[$i+1]
done

popd &> /dev/null

