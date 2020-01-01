#!/bin/bash
ZIP=$0

echo "Decompressing File"
sleep 1
mkdir dialogflow && tar -C -xvf dialogflow "$ZIP"

cd dialogflow/masters-go
bash deployments/setup.sh
