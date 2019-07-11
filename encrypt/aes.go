package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

// AESEncode 加密 AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
func AESEncode(key, plainText []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	plainText = pkcs5Padding(plainText, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	cipherData := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherData, plainText)
	return hex.EncodeToString(cipherData), nil
}

// AESDecode 解密
func AESDecode(key []byte, cipherText string) (string, error) {
	cipherData, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cipherData))
	blockMode.CryptBlocks(origData, cipherData)
	origData = pkcs5UnPadding(origData)
	return bytesToStr(origData), nil
}

// 填充明文
func pkcs5Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plainText, padtext...)
}

// 去除填充
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
