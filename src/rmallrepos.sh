#!/bin/bash

i="0"
no=$1
while [ $i -lt $no ]
do
    rm -rf "$i"
    i=$[$i+1]
done
