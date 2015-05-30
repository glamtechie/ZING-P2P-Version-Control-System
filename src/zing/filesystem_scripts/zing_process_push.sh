#!/bin/bash

#1 == patchname
pushd . &> /dev/null
cd .zing/global

git checkout master;

git am $1
if [ "$?" -ne "0" ]; then
    popd &> /dev/null
    exit 1
fi

popd &> /dev/null



