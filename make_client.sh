#!/bin/bash

cat ~/.ssh/id_rsa.pub | ssh glamtechie@137.110.91.190 'cat >> ~/.ssh/authorized_keys'
mkdir $1; cd $1;
git remote add origin glamtechie@137.110.91.190:/Users/glamtechie/Documents/dum
