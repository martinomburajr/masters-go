#!/bin/bash

echo "SSH!"

gnome-terminal -- /bin/bash -c "ssh -i ~/.ssh/id_rsa   martin.omburajr@35.188.247.90"
gnome-terminal -- /bin/bash -c "ssh -i ~/.ssh/id_rsa   martin.omburajr@35.203.187.69"
gnome-terminal -- /bin/bash -c "ssh -i ~/.ssh/id_rsa   martin.omburajr@35.231.221.212"
gnome-terminal -- /bin/bash -c "ssh -i ~/.ssh/id_rsa   martin.omburajr@35.195.41.51"

wait

echo "Logged In!"