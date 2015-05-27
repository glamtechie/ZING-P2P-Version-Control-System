#!/bin/bash


#1 == branchname, 2 == patch name
pushd . &> /dev/null;

cd .zing/global;

temp_head_first=`git log --pretty=oneline | sed -n '1p' | awk '{print $1}'`
echo $temp_head_first
git checkout -b temp &> /dev/null;

git pull ../../ $1  &> /dev/null;

git format-patch $temp_head_first --stdout > $2 


git checkout master &> /dev/null
git branch -D temp &> /dev/null


popd &> /dev/null;
