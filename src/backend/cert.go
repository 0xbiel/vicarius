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
