package encrypt

import (
	"crypto/hmac"
	"crypto/md5"  // nolint:gosec
	"crypto/sha1" // nolint:gosec
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

// 从[]byte生成md5密文
func MD5ByBytes(data []byte, salt ...[]byte) (cipher string, err error) {
	var s []byte
	if len(salt) > 0 {
		s = salt[0]
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(md5.New, s)
	} else {
		h = md5.New() // nolint:gosec
	}
	if _, err = h.Write(data); err != nil {
		return
	}
	cipher = hex.EncodeToString(h.Sum(nil))
	return
}

// 从string生成md5密文
func MD5ByStr(data string, salt ...string) (cipher string, err error) {
	var s []byte
	if len(salt) > 0 {
		s = strToBytes(salt[0])
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(md5.New, s)
	} else {
		h = md5.New() // nolint:gosec
	}
	if _, err = h.Write(strToBytes(data)); err != nil {
		return
	}
	cipher = hex.EncodeToString(h.Sum(nil))
	return
}

// 从[]string生成md5密文
func MD5ByStrings(data []string, salt ...string) (string, error) {
	var s []byte
	if len(salt) > 0 {
		s = strToBytes(salt[0])
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(md5.New, s)
	} else {
		h = md5.New() // nolint:gosec
	}
	for k := range data {
		_, err := h.Write([]byte(data[k]))
		if err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// 根据[]byte生成sha1密文
func SHA1ByBytes(data []byte, salt ...[]byte) (cipher string, err error) {
	var s []byte
	if len(salt) > 0 {
		s = salt[0]
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(sha1.New, s)
	} else {
		h = sha1.New() // nolint:gosec
	}
	if _, err = h.Write(data); err != nil {
		return
	}
	cipher = hex.EncodeToString(h.Sum(nil))
	return
}

// 根据string生成sha1密文
func SHA1ByStr(data string, salt ...string) (cipher string, err error) {
	var s []byte
	if len(salt) > 0 {
		s = strToBytes(salt[0])
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(sha1.New, s)
	} else {
		h = sha1.New() // nolint:gosec
	}
	if _, err = h.Write(strToBytes(data)); err != nil {
		return
	}
	cipher = hex.EncodeToString(h.Sum(nil))
	return
}

// 根据[]byte生成sha256密文
func SHA256ByBytes(data []byte, salt ...[]byte) (cipher string, err error) {
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
	if _, err = h.Write(data); err != nil {
		return
	}
	// 计算出字符串格式的签名
	cipher = hex.EncodeToString(h.Sum(nil))
	return
}

// 根据string生成sha256密文
func SHA256ByStr(data string, salt ...string) (cipher string, err error) {
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
	if _, err = h.Write(strToBytes(data)); err != nil {
		return
	}
	// 计算出字符串格式的签名
	cipher = hex.EncodeToString(h.Sum(nil))
	return
}
