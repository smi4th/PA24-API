#!/bin/bash

# this file is used to install the necessary packages for the project
# it is assumed that the user has sudo privileges
# it will also launch the API server

# !! IMPORTANT !!
# this `installation.sh` file will NOT launch docker containers

sudo apt update -y
sudo apt upgrade -y

# install go

# if there is no go folder
if [ ! -d "/usr/local/go" ]; then

    echo "Installing go..."

    sudo apt install -y wget
    wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
    rm go1.22.1.linux-amd64.tar.gz

    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

fi

cd PA24-API

echo "Launching API server..."
/usr/local/go/bin/go run . # run the API server