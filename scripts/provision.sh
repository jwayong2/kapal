#!/bin/bash

# Update, get python-software-properties in order to get add-apt-repository, 
# then update (for latest git version):
#sudo apt-get update
#sudo add-apt-repository -y ppa:git-core/ppa
sudo apt-get update

# Git, Golang, Docker & Curl:
sudo apt-get install -y git 
sudo apt-get install -y golang-go 
sudo apt-get install -y docker
sudo apt-get install -y curl

# For NFS speedup:
#sudo apt-get install -y nfs-common portmap

# Configure Go & Vim:
# http://tip.golang.org/misc/vim/readme.txt?m=text
# Configure Go workspace:
su vagrant -c "echo 'export GOROOT=/usr/lib/go' >> ~/.profile"
su vagrant -c "echo 'export GOPATH=/opt/gopath' >> ~/.profile"
su vagrant -c "echo 'export GOBIN=/usr/bin/go' >> ~/.profile"
su vagrant -c "echo 'export PATH=\$PATH:\$GOBIN' >> ~/.profile"

sudo chown -R vagrant:vagrant /opt/gopath

