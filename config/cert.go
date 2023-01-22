package config

import (
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"github.com/go-base-lib/coderutils"
)

//go:embed certs/ca.crt
var CaCertBytes []byte

//go:embed certs/client.crt
var ClientCertBytes []byte

//go:embed certs/client.key
var ClientKeyBytes []byte

//go:embed certs/client-server.key
var ClientServerKeyBytes string

//go:embed certs/proxy443.crt
var Proxy443CertBytes []byte

//go:embed certs/proxy443.key
var Proxy443KeyBytes []byte

var (
	ClientPublicKey  *rsa.PublicKey
	ClientPrivateKey *rsa.PrivateKey

	ClientServerPrivateKey *rsa.PrivateKey
)

func init() {
	var err error

	block, _ := pem.Decode(ClientCertBytes)
	if block == nil {
		panic("解析客户端证书失败")
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("解析证书信息失败: " + err.Error())
	}

	ClientPublicKey = certificate.PublicKey.(*rsa.PublicKey)

	ClientPrivateKey, err = coderutils.PEM2RsaPrivateKey(string(ClientKeyBytes))
	if err != nil {
		panic(err)
	}

	ClientServerPrivateKey, err = coderutils.PEM2RsaPrivateKey(ClientServerKeyBytes)
	if err != nil {
		panic(err)
	}
}
