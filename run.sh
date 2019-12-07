#!/usr/bin/env bash

# Check if directory exists
if [ -d "./data" ]
then
  echo "Deleting Data Directory"
  rm -rf ./data
else
  echo "Directory data does NOT exists."
fi

# Run Go Application
go run main.go

exit