#!/bin/bash

echo "Exporting!"

go build -a -v -o masters
sudo chmod 0777 masters
#
scp -i ~/.ssh/id_rsa masters   martin.omburajr@35.188.247.90:~/trial0/masters &
scp -i ~/.ssh/id_rsa forever.sh   martin.omburajr@35.188.247.90:~/trial0/ &
scp -i ~/.ssh/id_rsa clean.sh   martin.omburajr@35.188.247.90:~/trial0/ &
scp -i ~/.ssh/id_rsa -r _p1   martin.omburajr@35.188.247.90:~/trial0/ &

scp -i ~/.ssh/id_rsa masters   martin.omburajr@35.203.187.69:~/trial0/masters &
scp -i ~/.ssh/id_rsa forever.sh   martin.omburajr@35.203.187.69:~/trial0/ &
scp -i ~/.ssh/id_rsa clean.sh   martin.omburajr@35.203.187.69:~/trial0/ &
scp -i ~/.ssh/id_rsa -r _p2   martin.omburajr@35.203.187.69:~/trial0/ &
#
scp -i ~/.ssh/id_rsa masters   martin.omburajr@35.231.221.212:~/trial0/masters &
scp -i ~/.ssh/id_rsa forever.sh   martin.omburajr@35.231.221.212:~/trial0/ &
scp -i ~/.ssh/id_rsa clean.sh   martin.omburajr@35.231.221.212:~/trial0/ &
scp -i ~/.ssh/id_rsa -r _p3   martin.omburajr@35.231.221.212:~/trial0/ &
#
scp -i ~/.ssh/id_rsa masters   martin.omburajr@35.195.41.51:~/trial0/masters &
scp -i ~/.ssh/id_rsa forever.sh   martin.omburajr@35.195.41.51:~/trial0/ &
scp -i ~/.ssh/id_rsa clean.sh   martin.omburajr@35.195.41.51:~/trial0/ &
scp -i ~/.ssh/id_rsa -r _p4   martin.omburajr@35.195.41.51:~/trial0/ &
wait

echo "Exported!"