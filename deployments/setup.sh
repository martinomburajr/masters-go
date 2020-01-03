#!/bin/bash

sudo apt update && sudo apt install -y nano wget curl git r-base

wget https://dl.google.com/go/go1.13.5.linux-amd64.tar.gz;
sudo tar -xvf go1.13.5.linux-amd64.tar.gz
sudo mv go /usr/local
go get -d
export PATH=$PATH:/usr/local/go/bin

sudo Rscript -e 'install.packages("ggplot2",repos = "http://cran.us.r-project.org")'
sudo Rscript -e 'install.packages("readr",repos = "http://cran.us.r-project.org")'
sudo Rscript -e 'install.packages("knitr",repos = "http://cran.us.r-project.org")'
sudo Rscript -e 'install.packages("dplyr",repos = "http://cran.us.r-project.org")'

logout