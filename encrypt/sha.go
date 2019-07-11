package encrypt

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

// MD5ByBytes 从[]byte生成md5密文
func MD5ByBytes(value []byte, salt ...[]byte) (result string, err error) {
	var s []byte
	if len(salt) > 0 {
		s = salt[0]
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(md5.New, s)
	} else {
		h = md5.New()
	}
	if _, err = h.Write(value); err != nil {
		return
	}
	result = hex.EncodeToString(h.Sum(nil))
	return
}

// MD5ByStr 从string生成md5密文
func MD5ByStr(value string, salt ...string) (result string, err error) {
	var s []byte
	if len(salt) > 0 {
		s = strToBytes(salt[0])
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(md5.New, s)
	} else {
		h = md5.New()
	}
	if _, err = h.Write(strToBytes(value)); err != nil {
		return
	}
	result = hex.EncodeToString(h.Sum(nil))
	return
}

// SHA1ByBytes 根据[]byte生成sha1密文
func SHA1ByBytes(value []byte, salt ...[]byte) (result string, err error) {
	var s []byte
	if len(salt) > 0 {
		s = salt[0]
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(sha1.New, s)
	} else {
		h = sha1.New()
	}
	if _, err = h.Write(value); err != nil {
		return
	}
	result = hex.EncodeToString(h.Sum(nil))
	return
}

// SHA1ByStr 根据string生成sha1密文
func SHA1ByStr(value string, salt ...string) (result string, err error) {
	var s []byte
	if len(salt) > 0 {
		s = strToBytes(salt[0])
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(sha1.New, s)
	} else {
		h = sha1.New()
	}
	if _, err = h.Write(strToBytes(value)); err != nil {
		return
	}
	result = hex.EncodeToString(h.Sum(nil))
	return
}

// SHA256ByBytes 根据[]byte生成sha256密文
func SHA256ByBytes(value []byte, salt ...[]byte) (result string, err error) {
	var s []byte
	if len(salt) > 0 {
		s = salt[0]
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(sha256.New, s)
	} else {
		h = sha256.New()
	}
	if _, err = h.Write(value); err != nil {
		return
	}
	// 计算出字符串格式的签名
	result = hex.EncodeToString(h.Sum(nil))
	return
}

// SHA256ByStr 根据string生成sha256密文
func SHA256ByStr(value string, salt ...string) (result string, err error) {
	var s []byte
	if len(salt) > 0 {
		s = strToBytes(salt[0])
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(sha256.New, s)
	} else {
		h = sha256.New()
	}
	if _, err = h.Write(strToBytes(value)); err != nil {
		return
	}
	// 计算出字符串格式的签名
	result = hex.EncodeToString(h.Sum(nil))
	return
}
