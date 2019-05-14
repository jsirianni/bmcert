#!/usr/bin/env bash
cd /usr/local/bin
wget https://github.com/BlueMedoraPublic/bmcert/releases/download/0.2.0/bmcert-darwin-amd64.zip
unzip -o bmcert-darwin-amd64.zip
chmod +x bmcert
rm -f bmcert-darwin-amd64.zip
which bmcert
bmcert version
