#!/bin/bash


port=$1 #start portno
cloneip=$2
i=$3 #lower bound
no=$4 #upper bound
echo "no is"
echo $no

while [ $i -lt $no ]
do
    pushd . &> /dev/null
    mkdir "$i"
    cd "$i"
    zing clone $cloneip  $port
    zing_server $port  &
    sleep 5
    i=$[$i+1]
    port=$[$port+1]
    popd &> /dev/null
done


