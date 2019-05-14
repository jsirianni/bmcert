#!/usr/bin/env bash
cd /usr/local/bin
wget https://github.com/BlueMedoraPublic/bmcert/releases/download/v0.2.0/bmcert-v0.2.0-darwin-amd64.zip
unzip -o bmcert-v0.2.0-darwin-amd64.zip
chmod +x bmcert
rm -f bmcert-v0.2.0-darwin-amd64.zip
which bmcert
bmcert version
