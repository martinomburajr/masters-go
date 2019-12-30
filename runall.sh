#! /bin/bash

echo "Initializing"
sleep 1

params="_params"

for i in {0..3} ; do
  echo ${i}
  path="${params}/${i}"
  echo $path
  go run main.go --params=$path --parallelism=true &
  echo "################################################"
done
