#! /bin/bash

echo "Initializing"
sleep 1

params="_params"

for i in {0..1} ; do
  path="${params}/${i}"
  echo $path
  launch go run main.go --params=$path &
  while kill -0 "$PROC_ID" >/dev/null 2>&1; do
    echo "PROCESS IS RUNNING"
  done
  echo "PROCESS TERMINATED"
  exit 0

done
