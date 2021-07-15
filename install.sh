#!/usr/bin/env bash

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     os=linux;;
    Darwin*)    os=darwin;;
    *)          os="UNKNOWN:${unameOut}"
esac

TAG=v1.1.5
ZIP_FILE="bmcert_${TAG}_${os}_amd64.zip"
wget "https://github.com/BlueMedoraPublic/bmcert/releases/download/${TAG}/${ZIP_FILE}"
unzip -o -j "${ZIP_FILE}" "bmcert" -d "/usr/local/bin"
chmod +x bmcert
rm -f "${ZIP_FILE}"
which bmcert
bmcert version
