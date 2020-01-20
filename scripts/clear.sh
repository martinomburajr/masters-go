#!/usr/bin/env bash

if [ -d "./data" ]
then
  echo "Deleting Data Directory"
  rm -rf ./data
else
  echo "Directory data does NOT exists."
fi

exit