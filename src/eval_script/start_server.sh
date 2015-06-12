#!/bin/bash


i=$1 #lower bound
no=$2 #upper bound
echo "no is"
echo $no

while [ $i -lt $no ]
do
    pushd . &> /dev/null
    cd "$i"
    zing_server  &
    i=$[$i+1]
    sleep 3
    popd &> /dev/null
done


