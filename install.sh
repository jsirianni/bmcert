#!/usr/bin/env bash
TAG=v1.1.1
cd /usr/local/bin
wget https://github.com/BlueMedoraPublic/bmcert/releases/download/$TAG/bmcert-$TAG-darwin-amd64.zip
unzip -o bmcert-$TAG-darwin-amd64.zip
chmod +x bmcert
rm -f bmcert-$TAG-darwin-amd64.zip
which bmcert
bmcert version
