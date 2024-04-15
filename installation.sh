#!/bin/bash

# this file is used to install the necessary packages for the project
# it is assumed that the user has sudo privileges
# it will also launch the API server

# !! IMPORTANT !!
# this `installation.sh` file will launch docker container

apt update -y
apt upgrade -y

# install go

# if there is no go folder
if [ ! -d "/usr/local/go" ]; then

    echo "Installing go..."

    apt install -y wget
    wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
    rm go1.22.1.linux-amd64.tar.gz

    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

fi

cd PA24-API

if [[ "$(docker images)" == *"api-container"* ]]; then
    echo "API container already exists..."
else
    echo "Building the API server container..."
    docker build -t api-container .
fi

echo "Running the API server container..."
docker run -p 80:80 -it --rm --name api-container api-container