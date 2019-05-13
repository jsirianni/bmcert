#!/bin/sh
if [ -z "$VAULT_ADDR" ]
then
    echo "Failed to read VAULT_ADDR"
    exit 1
fi

if [ -z "$VAULT_CERT_URL" ]
then
    echo "Failed to read VAULT_CERT_URL"
    exit 1
fi

VERSION=`cat cmd/version.go | grep "const VERSION" | cut -c 17- | tr -d '"'`
if [ -z "$VERSION" ]
then
      echo "Failed to get version from cmd/const.go"
      exit 1
fi

echo "Building bmcert ${VERSION}"

docker build \
    --no-cache \
    --build-arg version=${VERSION} \
    --build-arg token=${VAULT_GITHUB_TOKEN} \
    --build-arg addr=${VAULT_ADDR} \
    --build-arg url=${VAULT_CERT_URL} \
    -t bmcert:${VERSION} .

docker create -ti --name artifacts bmcert:${VERSION} bash && \
    docker cp artifacts:/bmcert-v${VERSION}-linux-amd64.zip bmcert-v${VERSION}-linux-amd64.zip && \
    docker cp artifacts:/bmcert-v${VERSION}-darwin-amd64.zip bmcert-v${VERSION}-darwin-amd64.zip && \
    docker cp artifacts:/bmcert-v${VERSION}.SHA256SUMS bmcert-v${VERSION}.SHA256SUMS

# cleanup
docker rm -fv artifacts &> /dev/null
