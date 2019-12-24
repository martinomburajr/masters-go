#!/bin/bash

sudo apt update

wget https://dl.google.com/go/go1.13.5.linux-amd64.tar.gz;
tar -C /usr/local -xzf go1.13.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

sudo apt install r-base
