# bmcert
CLI for generating signed certificates using Vault.


## Overview
`bmcert` uses the `cobra` project for parsing command line options. https://github.com/spf13/cobra
```
bmcert --help
```


## Installation
`bmcert` should be placed somewhere in your systems `path` and made executable.

`bmcert` relies on two environment variables
```
# your github token
VAULT_GITHUB_TOKEN=<token here>

# vault server url
VAULT_ADDR=https://vault.mynet.com:8200

# full URL to vault server, including the pki path
VAULT_CERT_URL=https://vault.mynet.com:8200/v1/<pki endpoint>
```

## Usage

### Create Certificate
Call the `create` command to generate a certificate.

#### Generate x509 PEM
`bmcert` will generate a single file that contains the full certificate chain
```
bmcert create --hostname <fqdn>
```

#### Generate x509 PEM, cert and key files
If individual certificate and private key files are desired, use the `--format` flag
```
bmcert create --hostname <fqdn> --format cert
```

#### Generate pkcs12
pkcs12 is supported when passing `p12` or `pkcs12` for `--format`.
Password (optional) will be used to secure the certificate.
```
bmcert create --hostname bob.bluemedora.localnet --format p12

bmcert create --hostname bob.bluemedora.localnet --format p12 --password medora
```


### Flags

#### Global flags
Global flags are used for any command:
```
-h, --help                help for bmcert
    --tls-skip-verify     Disable certificate verification when communicating with the Vault API (Defaults to false)
    --verbose             Enable verbose output --verbose
```

#### Create flags
When calling `bmcert create`:
```
    --alt-names string    The requested Subject Alternative Names, in a comma-delimited list
-F, --format string       The keyfile formant to output. [pem, p12] (default "pem")
-h, --help                help for create
-H, --hostname string     The fully qualified hostname.
    --ip-sans string      The requested IP Subject Alternative Names, in a comma-delimited list
-O, --output-dir string   The directory to output to. Defaults to working directory.
-P, --password string     The password to protect pkcs12 (p12) certificates (optional)
    --uri-sans string     The requested URI Subject Alternative Names, in a comma-delimited list. (ALTHA: Not tested)
```


## Building from Source
Ensure your `GOPATH` is set appropriately and then run:
```
go get github.com:BlueMedora/bmcert
cd $GOPATH/src/github.com/BlueMedora/bmcert
go get .
go build
```
