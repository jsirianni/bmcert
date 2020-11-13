package cert

import (
    "strconv"
    "crypto/x509"
    "encoding/pem"

    "github.com/pkg/errors"
)

func (s SignedCertificate) validate() error {
    if err := validatePrivateKey([]byte(s.PrivateKey)); err != nil {
        return err
    }
    if err := validateCertificate([]byte(s.Certificate)); err != nil {
        return err
    }
    if err := validateCA([]byte(s.IssuingCa)); err != nil {
        return err
    }
    return nil
}

func validatePrivateKey(key []byte) error {
    privPem, _ := pem.Decode(key)
    if privPem.Type != "RSA PRIVATE KEY" {
        return errors.New("Expected Vault to return an RSA private key, got: " + privPem.Type)
    }
    if _, err := x509.ParsePKCS1PrivateKey(privPem.Bytes); err != nil {
        return errors.Wrap(err, "failed to parse private key: " + string(key))
    }
    return nil
}

func validateCertificate(crt []byte) error {
    c, _ := pem.Decode(crt)
    if c.Type != "CERTIFICATE" {
        return errors.New("Expected Vault to return a certifcate, got: " + c.Type)
    }

    certifcate, err := x509.ParseCertificates(c.Bytes)
    if err != nil {
        return errors.Wrap(err, "failed to parse certifcate: " + string(crt))
    }

    if len(certifcate) > 1 {
        return errors.New("expected Vault to return a single certifcate, got: " + strconv.Itoa(len(certifcate)))
    }

    return nil
}

func validateCA(ca []byte) error {
    c, _ := pem.Decode(ca)
    if c.Type != "CERTIFICATE" {
        return errors.New("Expected Vault to return a certifcate authority, got: " + c.Type)
    }

    if _, err := x509.ParseCertificates(c.Bytes); err != nil {
        return errors.Wrap(err, "failed to parse certifcate authority: " + string(ca))
    }


    return nil
}
