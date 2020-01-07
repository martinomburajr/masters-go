#!/bin/bash

IP=$1
ZIP=$2

echo "Compressing Archive"
sleep 1
cd .. && tar --exclude=".git" --exclude=".idea" -cvzf $ZIP masters-go

sleep 1
echo "Copying To Remote: $IP"
sleep 1

scp -i ~/.ssh/id_rsa "$ZIP"  martin.omburajr@"$IP":~/

echo "Connecting To Remote: $IP"
sleep 1
ssh -i ~/.ssh/id_rsa "martin.omburajr@$IP":~/
