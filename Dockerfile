
# staging environment retrieves dependencies and compiles
# for Linux and MacOS
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
    go get github.com/hashicorp/vault/sdk/helper/certutil

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
ENV VAULT_CERT_URL=$url

COPY --from=stage /build/src/bmcert/bmcert-linux bmcert

RUN apt-get update && apt-get install -y openssl

RUN ./bmcert create --hostname test.bluemedora.localnet --tls-skip-verify
RUN ./bmcert create --hostname test.bluemedora.localnet --tls-skip-verify --format p12
RUN ./bmcert create --hostname test.bluemedora.localnet --tls-skip-verify --format cert

RUN openssl x509 -in test.bluemedora.localnet.pem -text -noout

# build the release with an image that includes zip and sha256sum
FROM debian:stable

WORKDIR /
ARG version

RUN apt-get update && apt-get install zip -y

COPY --from=stage /build/src/bmcert/bmcert-linux bmcert
RUN zip bmcert-v${version}-linux-amd64.zip bmcert

COPY --from=stage /build/src/bmcert/bmcert-darwin bmcert
RUN zip bmcert-v${version}-darwin-amd64.zip bmcert

RUN sha256sum bmcert-v${version}-linux-amd64.zip >> bmcert-v${version}.SHA256SUMS
RUN sha256sum bmcert-v${version}-darwin-amd64.zip >> bmcert-v${version}.SHA256SUMS
