package encrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	publicKeyPrefix = "-----BEGIN PUBLIC KEY-----"
	publicKeySuffix = "-----END PUBLIC KEY-----"

	PKCS1Prefix = "-----BEGIN RSA PRIVATE KEY-----"
	PKCS1Suffix = "-----END RSA PRIVATE KEY-----"

	PKCS8Prefix = "-----BEGIN PRIVATE KEY-----"
	PKCS8Suffix = "-----END PRIVATE KEY-----"
)

// 生成RSA私钥
func GenerateRSAPrivateKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// 验证RSA私钥
func ValidateRSAPrivateKey(key *rsa.PrivateKey) error {
	if key == nil {
		return errors.New("无效的私钥")
	}
	return key.Validate()
}

// RSA私钥转为Base64
func RSAPrivateKeyToBase64(privateKey *rsa.PrivateKey, version uint8) (string, error) {
	var (
		keyBytes []byte
		err      error
	)
	err = ValidateRSAPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	switch version {
	case 1:
		keyBytes = x509.MarshalPKCS1PrivateKey(privateKey)
	case 8:
		keyBytes, err = x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("仅支持转为PKCS的1和8版本的密钥")
	}
	return base64.RawURLEncoding.EncodeToString(keyBytes), nil
}

// RSA私钥转为Hex
func RSAPrivateKeyToHex(privateKey *rsa.PrivateKey, version uint8) (string, error) {
	var (
		keyBytes []byte
		err      error
	)
	switch version {
	case 1:
		keyBytes = x509.MarshalPKCS1PrivateKey(privateKey)
	case 8:
		keyBytes, err = x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("仅支持转为PKCS的1和8版本的密钥")
	}
	return hex.EncodeToString(keyBytes), nil
}

// 从RSA私钥中获取公钥并将公钥转为Base64字符串
func RSAPublicKeyToBase64(privateKey *rsa.PrivateKey) (string, error) {
	err := ValidateRSAPrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	keyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	return base64.RawURLEncoding.EncodeToString(keyBytes), nil
}

// 从RSA私钥中获取公钥并将公钥转为Hex字符串
func RSAPublicKeyToHex(privateKey *rsa.PrivateKey) (string, error) {
	err := ValidateRSAPrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	keyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	return hex.EncodeToString(keyBytes), nil
}

// 从base64解析RSA私钥
func Base64ToRSAPrivateKey(base64Str string) (privateKey *rsa.PrivateKey, version uint8, err error) {
	var (
		keyBytes []byte
	)
	keyBytes, err = base64.RawURLEncoding.DecodeString(base64Str)
	if err != nil {
		return
	}

	privateKey, version, err = ParseRSAPrivateKey(keyBytes)
	if err != nil {
		return
	}

	return
}

// Base64转为RSA公钥
func Base64ToRSAPublicKey(base64Str string) (publicKey *rsa.PublicKey, err error) {
	var keyBytes []byte
	keyBytes, err = base64.RawURLEncoding.DecodeString(base64Str)
	if err != nil {
		return
	}

	publicKey, err = x509.ParsePKCS1PublicKey(keyBytes)
	return
}

// Hex转为RSA私钥
func HexToRSAPrivateKey(hexStr string) (privateKey *rsa.PrivateKey, version uint8, err error) {
	var keyBytes []byte
	keyBytes, err = hex.DecodeString(hexStr)
	if err != nil {
		return
	}

	privateKey, version, err = ParseRSAPrivateKey(keyBytes)
	return
}

// Hex转为RSA公钥
func HexToRSAPublicKey(hexStr string) (publicKey *rsa.PublicKey, err error) {
	var keyBytes []byte
	keyBytes, err = hex.DecodeString(hexStr)
	if err != nil {
		return
	}

	publicKey, err = x509.ParsePKCS1PublicKey(keyBytes)
	return
}

// 解析RSA公钥文件
func ParseRSAPublicKeyFile(filePath string) (publicKey *rsa.PublicKey, err error) {
	var (
		file   []byte
		blocks []*pem.Block
	)
	file, err = ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return
	}

	blocks = ParsePEMBlocks(file)

	publicKey, err = ParseRSAPublicKey(blocks[0].Bytes)
	return
}

// 解析RSA私钥文件
func ParseRSAPrivateKeyFile(filePath string) (privateKey *rsa.PrivateKey, version uint8, err error) {
	var (
		file   []byte
		blocks []*pem.Block
	)
	file, err = ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return
	}

	blocks = ParsePEMBlocks(file)

	privateKey, version, err = ParseRSAPrivateKey(blocks[0].Bytes)
	return
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
func ParseRSAPublicKey(data []byte) (publicKey *rsa.PublicKey, err error) {
	var (
		parsedKey interface{}
		ok        bool
		cert      *x509.Certificate
	)

	parsedKey, err = x509.ParsePKIXPublicKey(data)
	if err != nil {
		cert, err = x509.ParseCertificate(data)
		if err == nil {
			parsedKey = cert.PublicKey
		} else {
			return
		}
	}

	publicKey, ok = parsedKey.(*rsa.PublicKey)
	if !ok {
		err = errors.New("不是有效的RSA公钥")
		return
	}

	return
}

// 解析PKCS1私钥
func ParsePKCS1PrivateKey(data []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(data)
}

// 解析PKCS8私钥
func ParsePKCS8PrivateKey(data []byte) (privateKey *rsa.PrivateKey, err error) {
	var (
		parsedKey interface{}
		ok        bool
	)
	parsedKey, err = x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		return
	}

	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		err = errors.New("不是有效的RSA私钥")
		return
	}

	return
}

// 解析RSA私钥，自动识别PKCS1和PKCS8
func ParseRSAPrivateKey(data []byte) (privateKey *rsa.PrivateKey, version uint8, err error) {
	var (
		pkcs8key interface{}
		ok       bool
	)

	version = 1
	// 尝试PKCS1
	privateKey, err = x509.ParsePKCS1PrivateKey(data)
	if err != nil {
		// 如果PKCS1和PKCS8，直认为无效
		if !strings.Contains(err.Error(), "PKCS8") {
			return
		}
		// 尝试PKCS8
		pkcs8key, err = x509.ParsePKCS8PrivateKey(data)
		if err != nil {
			return
		}
		version = 8
		// 断言为私钥
		if privateKey, ok = pkcs8key.(*rsa.PrivateKey); !ok {
			err = errors.New("不是有效的RSA私钥")
			return
		}
	}
	return
}

// 格式化RSA公钥
func FormatRSAPublicKey(key string) []byte {
	return formatKey(key, publicKeyPrefix, publicKeySuffix, 64)
}

// 格式化PKCS1私钥
func FormatPKCS1PrivateKey(key string) []byte {
	return formatKey(key, PKCS1Prefix, PKCS1Suffix, 64)
}

// 格式化PKCS8私钥
func FormatPKCS8PrivateKey(key string) []byte {
	return formatKey(key, PKCS8Prefix, PKCS8Suffix, 64)
}

func formatKey(raw, prefix, suffix string, lineCount int) []byte {
	var err error
	raw = strings.Replace(raw, PKCS8Prefix, "", 1)
	raw = strings.Replace(raw, PKCS8Suffix, "", 1)
	if raw == "" {
		return nil
	}
	raw = strings.Replace(raw, prefix, "", 1)
	raw = strings.Replace(raw, suffix, "", 1)
	raw = strings.ReplaceAll(raw, " ", "")
	raw = strings.ReplaceAll(raw, "\n", "")
	raw = strings.ReplaceAll(raw, "\r", "")
	raw = strings.ReplaceAll(raw, "\t", "")

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
