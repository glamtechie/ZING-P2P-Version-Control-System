#!/bin/bash

pushd . &> /dev/null;

cd .zing/global;

git checkout -b temp &> /dev/null;
git pull ../../ $1 &> /dev/null;

popd &> /dev/null;


