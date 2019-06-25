
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
ENV GOPATH=/build

ADD . /build/src/bmcert

RUN \
    go get github.com/spf13/cobra && \
    go get github.com/BlueMedoraPublic/go-pkcs12 && \
    go get github.com/hashicorp/vault/sdk/helper/certutil && \
    go get github.com/mitchellh/go-homedir

RUN go test ./...

RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bmcert-linux
RUN env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bmcert-darwin


# perform tests with compiled binary
FROM debian:stable AS test

ARG token
ARG addr
ARG url
ENV VAULT_GITHUB_TOKEN=$token
ENV VAULT_ADDR=$addr
ENV VAULT_SKIP_VERIFY=true
ENV VAULT_CERT_URL=$url

COPY --from=stage /build/src/bmcert/bmcert-linux bmcert

RUN apt-get update >> /dev/null && \
    apt-get install -y openssl wget unzip >> /dev/null

RUN \
    wget -q https://releases.hashicorp.com/vault/1.1.2/vault_1.1.2_linux_amd64.zip && \
    unzip vault_1.1.2_linux_amd64.zip && mv vault /usr/bin && chmod +x /usr/bin/vault

RUN \
    vault login -method=github token=$VAULT_GITHUB_TOKEN >> /dev/null

# create and validate
RUN ./bmcert create --hostname test.bluemedora.localnet --tls-skip-verify && \
    openssl x509 -in test.bluemedora.localnet.pem -text -noout

RUN ./bmcert create --hostname test.bluemedora.localnet --tls-skip-verify --format cert && \
    openssl x509 -in test.bluemedora.localnet.crt -text -noout && \
    openssl rsa -in test.bluemedora.localnet.key -check

RUN ./bmcert create --hostname test.bluemedora.localnet --tls-skip-verify --format p12 --password password

RUN ./bmcert create -f --hostname test.bluemedora.localnet --tls-skip-verify --format p12 --password password

# build the release with an image that includes zip and sha256sum
FROM debian:stable

WORKDIR /
ARG version

RUN apt-get update >> /dev/null && \
    apt-get install zip -y >> /dev/null

COPY --from=stage /build/src/bmcert/bmcert-linux bmcert
RUN zip bmcert-v${version}-linux-amd64.zip bmcert

COPY --from=stage /build/src/bmcert/bmcert-darwin bmcert
RUN zip bmcert-v${version}-darwin-amd64.zip bmcert

RUN sha256sum bmcert-v${version}-linux-amd64.zip >> bmcert-v${version}.SHA256SUMS
RUN sha256sum bmcert-v${version}-darwin-amd64.zip >> bmcert-v${version}.SHA256SUMS
