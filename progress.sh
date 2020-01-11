#! /bin/bash

red=$'\e[1;31m' # Error
grn=$'\e[1;32m' # Success
yel=$'\e[1;33m' # Warnings
blu=$'\e[1;34m' # Info
mag=$'\e[1;35m' # Title
cyn=$'\e[1;36m'
end=$'\e[0m'

printf "${cyn}\n\
##########################################################\n \
                Showing Simulation Progress!\n\
##########################################################\n
${end}\n"

printf "${blu}Building Go Binary: ...\n${end}"
sleep 1
rm -rf ./masters-go
go build -a -v -o masters-go-progress
wait
sleep 3

printf "${yel}\n\Showing Progress: ...\n${end}"
sleep 3
for i in {0..1000} ; do
  ./masters-go-progress --params="_params" --showProgress=true --parallelism=true --dataDir="data" --logging=true --runstats=true &
  sleep 5
done