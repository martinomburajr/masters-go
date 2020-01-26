#!/bin/bash

echo "Exporting!"

go build -a -v -o masters
scp -i ~/.ssh/id_rsa masters   martin.omburajr@34.67.243.122:~/trial2-ext/masters &
scp -i ~/.ssh/id_rsa masters   martin.omburajr@35.231.221.212:~/trial2-ext/masters &
wait

echo "Exported!"