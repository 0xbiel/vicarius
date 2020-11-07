package cert

import (
  "bytes"
  "crypto"
  "crypto/rand"
  "crypto/rsa"
  "crypto/sha1"
  "crypto/tls"
  "crypto/x509"
  "crypto/x509/pkix"
  "encoding/pem"
  "errors"
  "fmt"
  "math/big"
  "net"
  "os"
  "path/filepath"
  "time"
)

// create serial number for the certificate.
var serialNumber = big.NewInt(0).SetBytes(bytes.Repeat([]byte{255}, 20))

type CertConfig struct {
  privKey	  rsa.PrivateKey
  keyId		  []byte
  privCert	  crypto.PrivateKey
  cert		  *x509.Certificate
}

func certCfg(cert *x509.Certificate, certPrivKey crypto.PrivateKey) (*CertConfig, error) {

  priv, err := rsa.GenerateKey(rand.Reader, 2048)
  if(err != nil) {
	return nil, err
  }

  pub := priv.Public()

  pkixKey, err := x509.MarshalPKIXPublicKey(pub)

  if(err != nil) {
	return nil, err
  }

  hash := sha1.New()
  hash.Write(pkixKey)
  keyId := hash.Sum(nil)

  return &CertConfig {
	privKey:  priv,
	keyId:	  keyId,
	privCert: certPrivKey,
	cert:	  cert,
  }, nil
}

// @@@TODO: lcCert function.
