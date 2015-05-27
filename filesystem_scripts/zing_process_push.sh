#!/bin/bash

#1 == patchname
pushd . &> /dev/null
cd .zing/global

git checkout master;

git am $1

popd &> /dev/null



