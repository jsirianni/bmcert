# bmcert
CLI for generating signed certificates using Vault.


## Overview
bmcert uses the `cobra` project for parsing command line options. https://github.com/spf13/cobra
```
bmcert --help
```


## Installation
bmcert should be placed somewhere in your systems `path` and made executable. bmcert relies
on an environment variable `VAULT_TOKEN`, which must have permission to create certificates.


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
--help                help for bmcert
--pkipath string      The vault certificate authority mount point (default "/v1/bm-pki-int/issue/bluemedora-dot-localnet")
--tls                 Enable or disable TLS encryption "--tls=true" (default true)
--vault-host string   The vault server (default "vault.bluemedora.localnet")
--vault-port string   The vault http port (default "8200")
--verbose             Enable verbose output --verbose
```

#### Create flags
When calling `bmcert create`:
```
--alt-names string    The requested Subject Alternative Names, in a comma-delimited list
--format string       The keyfile formant to output. [pem, p12] (default "pem")
--help                help for create
--hostname string     The fully qualified hostname.
--ip-sans string      The requested IP Subject Alternative Names, in a comma-delimited list
--output-dir string   The directory to output to. Defaults to working directory.
--uri-sans string     The requested URI Subject Alternative Names, in a comma-delimited list. (ALTHA: Not tested)
```


## Building from Source
```
git@github.com:BlueMedora/bmcert.git
cd bmcert
go get
go build
```
