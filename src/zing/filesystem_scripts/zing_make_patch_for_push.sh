#!/bin/bash


#1 == branchname, 2 == patch name
pushd . &> /dev/null;

cd .zing/global;

temp_head_first=`git log --pretty=oneline | sed -n '1p' | awk '{print $1}'`
echo $temp_head_first
git checkout -b temp &> /dev/null;
if [ "$?" -ne "0" ]; then
    exit 1
fi


git pull ../../ $1  &> /dev/null;

if [ "$?" -ne "0" ]; then
    exit 1
fi

git format-patch $temp_head_first --stdout > $2 
if [ "$?" -ne "0" ]; then
    exit 1
fi



git checkout master &> /dev/null
if [ "$?" -ne "0" ]; then
    exit 1
fi


git branch -D temp &> /dev/null
if [ "$?" -ne "0" ]; then
    exit 1
fi



popd &> /dev/null;

exit $?
