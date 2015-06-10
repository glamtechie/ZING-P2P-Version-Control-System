#!/bin/bash

port=27321
no=$1
echo "no is"
echo $no
i="0"

while [ $i -lt $no ]
do
    echo $port >> dummy
    i=$[$i+1]
    port=$[$port+1]
done
