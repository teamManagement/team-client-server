package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/go-base-lib/coderutils"
	"testing"
)

func TestDecryptClientServerData(t *testing.T) {
	encData := "hkQvwguBDmeRhqlibSyLN8uEo9y+TwLed2L3QVR3ev5Z6Af6m3q2UB5rJ7jfcYlYqQzz4MgH+orL8fmGDjjQY8Ln/lXVOLpP5P1qlOjCP3Kkh40NK777mpCaAoJ6sZB6GJdnpWWwPmer7ZkO/QjjyG7lUqFMEf44d+0epD0eoQD4s27M0zNYRAXxYurdnB8zZrREWD0xSnqz4thtxL9DF78FIAQaGL6Znp+nxFgsKjJoRy4zgWgRpWmLTuqlklt5l3+JiWe2poRODWe4EPtJzXuImeAoyT8LqqsVVg/bAqefquDP85k4fMenYJoB8A8kU8VipPVqwSM6Cy/phD1DQORqfOSKh9/oEpQvNm7G2X2H03MQDn19h/ehtAqNGPJaS00sLM6hye/2g0vQ1NCUHO7QJHw2r35tjN1u426efea7ztw6C+2PZ2BZ51LxPBB74/MjG7FqkmarpuBjD9t85WXl6pJZxnLhLjq2pSvsIzahAMRnOqliG2a00oH3T334dzl7DPZuZiS1ADG+aSY+BU5DvzjKzA24rDYUHy9um0qZRTTic/bnXyPwtBRcZfpzUIBnz48pRt6m3HkJI6laMvryuRKeb0Z5LZCGnuQpEt9Bxnu88z7/SBeLChdDXxd18ltemMyei+Ata3GbIkDOX+KWHZXlrJd0Jina6hizrIE="

	block, _ := pem.Decode([]byte(ClientServerKeyBytes))
	if block == nil {
		panic("解析私钥格式失败")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	decodeString, err := base64.StdEncoding.DecodeString(encData)
	if err != nil {
		panic(err)
	}

	data, err := rsa.DecryptPKCS1v15(rand.Reader, key, decodeString)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

}

func TestEncrypt(t *testing.T) {

	key := coderutils.Sm4RandomKey()

	encrypt, err := coderutils.RsaEncrypt(key, ClientPublicKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(encrypt))

}
