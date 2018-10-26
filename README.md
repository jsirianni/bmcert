# bmcert
CLI for generating signed certificates using Vault.

## Overview
bmcert uses the `cobra` project for parsing command line options. https://github.com/spf13/viper
```
bmcert --help
```

## Installation
bmcert should be placed somewhere in your systems `path` and made executable. bmcert relies
on an environment variable `VAULT_TOKEN`, which must have permission to create certificates.

## Usage
Call the `create` command to generate a certificate.


### Generate x509 PEM
short hostname, output to working directory
```
bmcert create --hostname myhost
```
fqdn, output to working directory
```
bmcert create --hostname myhost.bluemedora.localnet
```
short hostname, output to /etc/certs (both are valid paths)
```
bmcert create --hostname myhost --output-dir /etc/cert
bmcert create --hostname myhost --output-dir /etc/cert/
```
verbose output
```
bmcert create --hostname myhost --verbose
```

## Building from Source
```
git@github.com:BlueMedora/bmcert.git
cd bmcert
go get
go build
```
