#!/bin/bash

echo "Pulling changes to $1"
git pull .zing/global $1 | sed "s/git/zing/g";



