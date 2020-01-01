#! /bin/bash

echo "Initializing"
sleep 1

params="_params"


for i in {0..3} ; do
  echo ${i}
  path="${params}/${i}"
  echo $path
  /usr/local/go/bin/go run main.go --params=$path --parallelism=true --folder=${i} &
  echo "################################################"
done
