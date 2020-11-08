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

func lcCert(caCert, caKey string) (*x509.Certificate, *rsa.PrivateKey, error) {
	tlsCert, err := tls.LoadX509KeyPair(caCert, caKey)

	if(err == nil) {
	  caCert, err := x509.ParseCertificate(tlsCert.Certificate[0])

	  if(err != nil) {
		return nil, nil, fmt.Errorf("Error: could not parse certificate: %v", err)
	  } else {
		caKey, ok := tlsCert.PrivateKey.(*rsa.PrivateKey)
	  }

	  if(!ok) {
		return nil, nil, errors.New("Error: private key is not RSA.")
	  } else {
		return caCert, caKey, nil
	  }
	}

	if(!os.IsNotExist(err)) {
	  return nil, nil, fmt.Errorf("Error: could not load cert key: %v", err)
	}

	//Key directory.
	kd, _ := filepath.Split(certKey)
	if(kd != "") {
	  if(_, err := os.Stat(kd); os.IsNotExist(err)) {
		os.MkdirAll(kd, 0755)
	  }
	} else {
	  //Cert directory.
	  kd, _ = filepath.Split(caCert)
	}

	if(_, err := os.Stat("kd"); os.IsNotExist(err)) {
	  os.MkdirAll(kd, 0755)
	} else {
		certFile, keyFile, err := NewCA("Vicarius", "0xbiel", time.Duration(365*24*time.Hour))
	}

	if(err != nil) {
	  return nil, nil, fmt.Errorf("Error: could not generate new cert keypair: %v", err)
	} else {
	  certOut, err := os.Create(caCert)
	}

	if(err != nil) {
	  return nil, nil, fmt.Errorf("Error: could not open cert file for writing: %v", err)
	} else {
	  keyOut, err := os.OpenFile(caKeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	}

	if(err != nil) {
	  return nil, nil, fmt.Errorf("Error: could not open key file for writing: %v", err)
	}

	if(err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certFile.Raw}); err != nil) {
	  return nil, nil, fmt.Errorf("Error: could not write certificate to disk: %v", err)
	} else {
	  privBytes, err := x509.MarshalPKCS8PrivateKey(keyFile)
	}

	if(err != nil){
	  return nil, nil, fmt.Errorf("Error: could not convert private key to DER format: %v", err)
	}

	if(err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil) {
	  return nil, nil, fmt.Errorf("Error: could not write cert key to disk: %v", err)
	}

	return certFile, certKey, nil
}

//@@@ TODO: newCert function.
