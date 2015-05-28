#!/bin/bash

#1 == patchname
pushd . &> /dev/null

cd .zing/global
git checkout master;

git pull ../../ $1;


popd &> /dev/null



