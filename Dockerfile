
# staging environment compiles for Linux and MacOS
FROM golang:1.12 AS stage

WORKDIR /build/src/bmcert
ARG token
ARG addr
ARG url
ARG version
ENV VAULT_GITHUB_TOKEN=$token
ENV VAULT_ADDR=$addr
ENV VAULT_CERT_URL=$url
ENV VAULT_SKIP_VERIFY=true
ENV GOPATH=/build

ADD . /build/src/bmcert

RUN \
    go get github.com/spf13/cobra && \
    go get github.com/BlueMedoraPublic/go-pkcs12 && \
    go get github.com/hashicorp/vault/sdk/helper/certutil && \
    go get github.com/mitchellh/go-homedir

RUN go test ./...

# build without cgo, we do not need it for bmcert
RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bmcert
RUN env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bmcert-darwin

RUN apt-get update >> /dev/null && \
    apt-get install -y openssl wget unzip zip >> /dev/null

RUN \
    wget -q https://releases.hashicorp.com/vault/1.1.2/vault_1.1.2_linux_amd64.zip && \
    unzip vault_1.1.2_linux_amd64.zip && mv vault /usr/bin && chmod +x /usr/bin/vault

RUN \
    vault login -method=github token=$VAULT_GITHUB_TOKEN >> /dev/null

# create and validate
RUN ./bmcert create --hostname test3.bluemedora.localnet --tls-skip-verify && \
    openssl x509 -in test3.bluemedora.localnet.pem -text -noout

RUN ./bmcert create --hostname test2.bluemedora.localnet --tls-skip-verify --format cert && \
    openssl x509 -in test2.bluemedora.localnet.crt -text -noout && \
    openssl rsa -in test2.bluemedora.localnet.key -check

RUN ./bmcert create --hostname test2.bluemedora.localnet --tls-skip-verify --format p12 --password password

# test force replace flag
RUN ./bmcert create -f --hostname test2.bluemedora.localnet --tls-skip-verify --format p12 --password password

# test cert expiration, current year and future year should not be equal
RUN \
    ./bmcert create --hostname test2.bluemedora.localnet --tls-skip-verify -f --ttl 12m && \
    CURRENT_YEAR=$(TZ=GMT date +"%c %Z" | awk '{print $5}') && \
    FUTURE_YEAR=$(openssl x509 -in test2.bluemedora.localnet.pem -text -noout -dates | grep notAfter | awk '{print $4}') && \
    if [ "$CURRENT_YEAR" = "$FUTURE_YEAR" ]; then exit 1; fi

# test cert expiration, current year and future year should be equal
RUN \
    ./bmcert create --hostname test2.bluemedora.localnet --tls-skip-verify -f --ttl 1s && \
    CURRENT_YEAR=$(TZ=GMT date +"%c %Z" | awk '{print $5}') && \
    FUTURE_YEAR=$(openssl x509 -in test2.bluemedora.localnet.pem -text -noout -dates | grep notAfter | awk '{print $4}') && \
    if [ "$CURRENT_YEAR" != "$FUTURE_YEAR" ]; then exit 1; fi

# build the relese
#
RUN zip bmcert-v${version}-linux-amd64.zip bmcert
RUN mv bmcert-darwin bmcert && zip bmcert-v${version}-darwin-amd64.zip bmcert
RUN sha256sum bmcert-v${version}-linux-amd64.zip >> bmcert-v${version}.SHA256SUMS
RUN sha256sum bmcert-v${version}-darwin-amd64.zip >> bmcert-v${version}.SHA256SUMS
