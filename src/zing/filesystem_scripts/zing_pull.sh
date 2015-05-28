#!/bin/bash

echo "Pulling changes to $1"
git pull origin $1 | sed "s/git/zing/g";



