#! /bin/bash

#!/bin/bash
red=$'\e[1;31m' # Error
grn=$'\e[1;32m' # Success
yel=$'\e[1;33m' # Warnings
blu=$'\e[1;34m' # Info
mag=$'\e[1;35m' # Title
cyn=$'\e[1;36m'
end=$'\e[0m'

for j in {0..1} ; do
  printf "${mag}\n\
  ##########################################################\n \
                  Starting Simulation ${j} \n\
  ##########################################################\n
  ${end}\n"

  printf "${blu}\n\nInitializing Run: ...\n${end}"
  sleep 2
  for i in {0..0} ; do
    ./masters --params="_params" --numWorkers=2 --parallelism=true --dataDir="data" --logging=true
    printf "${yel}##################################### RUN: ${i} COMPLETE########################################\n${end}"
  done
done