#!/bin/bash

#Initialize local repository
git init &> /dev/null;

#Add .zing dir that will contain all metadata to gitignore
echo ".zing" > .gitignore
echo ".gitignore" >> .gitignore

if [ -d ".zing" ]; then
    rm -rf .zing
    echo "Reinitialized existing zing repository in $PWD/.zing"
else
    echo "Initialized zing repository in $PWD/.zing"
fi

pushd . &> /dev/null;

mkdir .zing && cd .zing;

#This ID needs to be a unique per node identifier
echo "ID:#" > config

#global dir that hold the pushes
mkdir global && cd global;
git init &> /dev/null;
git commit --allow-empty --allow-empty-message -m '' &> /dev/null
git tag -a -m '' ROOT
#git config receive.denyCurrentBranch ignore
popd &> /dev/null

git remote add origin .zing/global/
git pull origin master
