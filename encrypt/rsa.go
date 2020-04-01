package encrypt

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	kPublicKeyPrefix = "-----BEGIN PUBLIC KEY-----"
	kPublicKeySuffix = "-----END PUBLIC KEY-----"

	kPKCS1Prefix = "-----BEGIN RSA PRIVATE KEY-----"
	KPKCS1Suffix = "-----END RSA PRIVATE KEY-----"

	kPKCS8Prefix = "-----BEGIN PRIVATE KEY-----"
	KPKCS8Suffix = "-----END PRIVATE KEY-----"
)

// 解析RSA Public Key 文件
func ParseRSAPublicKeyFile(filePath string) (*rsa.PublicKey, error) {
	file, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}

	block := ParsePEMBlocks(file)

	publicKey, err := ParseRSAPublicKey(block[0].Bytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

/*
 解析RSA Private Key 文件
 第二个出参是key的类型（PKCS1/PKCS8）
*/
func ParseRSAPrivateKeyFile(filePath string) (*rsa.PrivateKey, string, error) {
	file, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, "", err
	}

	block := ParsePEMBlocks(file)

	privateKey, keyType, err := ParseRSAPrivateKey(block[0].Bytes)
	if err != nil {
		return nil, "", err
	}
	return privateKey, keyType, nil
}

// ParsePEMBlocks 解析PEM区块
func ParsePEMBlocks(data []byte) []*pem.Block {
	var (
		blocks []*pem.Block
		block  *pem.Block
		rest   []byte
	)
	block, rest = pem.Decode(data)
	if block != nil {
		blocks = append(blocks, block)
		for len(rest) > 0 {
			block, rest = pem.Decode(rest)
			if block != nil {
				blocks = append(blocks, block)
			}
		}
	}
	return blocks
}

// // ParseCertificate 解析x509证书文件
// func ParseCertificate(data []byte) (*x509.Certificate, error) {
// 	cert, err := x509.ParseCertificate(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return cert, nil
// }

// ParseRSAPublicKey 解析RSA公钥
func ParseRSAPublicKey(data []byte) (*rsa.PublicKey, error) {
	var err error
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(data); err != nil {
		if cert, err := x509.ParseCertificate(data); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	var publicKey *rsa.PublicKey
	var ok bool
	if publicKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, errors.New("不是有效的RSA公钥")
	}

	return publicKey, nil
}

// 解析PKCS1私钥
func ParsePKCS1PrivateKey(data []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(data)
}

// 解析PKCS8私钥
func ParsePKCS8PrivateKey(data []byte) (*rsa.PrivateKey, error) {
	parsedKey, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		return nil, err
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("不是有效的RSA私钥")
	}

	return privateKey, nil
}

// 解析RSA私钥，自动识别PKCS1/PKCS8
func ParseRSAPrivateKey(data []byte) (*rsa.PrivateKey, string, error) {
	var err error
	var parsedKey interface{}
	var t string = ""

	// 尝试PKCS1
	parsedKey, err = x509.ParsePKCS1PrivateKey(data)
	if err != nil {
		if err.Error() != "x509: failed to parse private key (use ParsePKCS8PrivateKey instead for this key format)" {
			return nil, "", err
		}
		parsedKey, err = x509.ParsePKCS8PrivateKey(data)
		if err != nil {
			return nil, "", err
		}
		t = "PKCS8"
	}

	if t == "" {
		t = "PKCS1"
	}
	var privateKey *rsa.PrivateKey
	var ok bool
	if privateKey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, "", errors.New("不是有效的RSA私钥")
	}

	return privateKey, t, nil
}

// 格式化RSA公钥
func FormatRSAPublicKey(key string) []byte {
	return formatKey(key, kPublicKeyPrefix, kPublicKeySuffix, 64)
}

// 格式化PKCS1私钥
func FormatPKCS1PrivateKey(key string) []byte {
	return formatKey(key, kPKCS1Prefix, KPKCS1Suffix, 64)
}

// 格式化PKCS8私钥
func FormatPKCS8PrivateKey(key string) []byte {
	return formatKey(key, kPKCS8Prefix, KPKCS8Suffix, 64)
}

func formatKey(raw, prefix, suffix string, lineCount int) []byte {
	var err error
	raw = strings.Replace(raw, kPKCS8Prefix, "", 1)
	raw = strings.Replace(raw, KPKCS8Suffix, "", 1)
	if raw == "" {
		return nil
	}
	raw = strings.Replace(raw, prefix, "", 1)
	raw = strings.Replace(raw, suffix, "", 1)
	raw = strings.Replace(raw, " ", "", -1)
	raw = strings.Replace(raw, "\n", "", -1)
	raw = strings.Replace(raw, "\r", "", -1)
	raw = strings.Replace(raw, "\t", "", -1)

	var sl = len(raw)
	var c = sl / lineCount
	if sl%lineCount > 0 {
		c++
	}

	var buf bytes.Buffer
	if _, err = buf.WriteString(prefix + "\n"); err != nil {
		return nil
	}
	for i := 0; i < c; i++ {
		var b = i * lineCount
		var e = b + lineCount
		if e > sl {
			if _, err = buf.WriteString(raw[b:]); err != nil {
				return nil
			}
		} else {
			if _, err = buf.WriteString(raw[b:e]); err != nil {
				return nil
			}
		}
		if _, err = buf.WriteString("\n"); err != nil {
			return nil
		}
	}
	if _, err = buf.WriteString(suffix); err != nil {
		return nil
	}
	return buf.Bytes()
}
