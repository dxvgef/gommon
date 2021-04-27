package random

import (
	"bytes"
	cRand "crypto/rand"
	"errors"
	"math/big"
	"math/rand"
	"time"
	"unsafe"
)

// Upper 指定长度的随机大写字母
func Upper(l int) string {
	var (
		result bytes.Buffer
		temp   string
	)
	for i := 0; i < l; {
		ii, err := Int32(65, 90)
		if err != nil {
			return ""
		}
		if string(ii) != temp {
			temp = string(ii)
			if _, err := result.WriteString(temp); err != nil {
				return ""
			}
			i++
		}
	}
	return result.String()
}

// Lower 指定长度的随机小写字母
func Lower(l int) string {
	var (
		result bytes.Buffer
		temp   string
	)
	for i := 0; i < l; {
		ii, err := Int32(97, 122)
		if err != nil {
			return ""
		}
		if string(ii) != temp {
			temp = string(ii)
			if _, err := result.WriteString(temp); err != nil {
				return ""
			}
			i++
		}
	}
	return result.String()
}

// CustomString，指定长度的随机字符串，第二个参数限制只能出现指定的字符
func CustomString(l int, specifiedStr string) string {
	const (
		letterIdxBits = 6
		letterIdxMask = 1<<letterIdxBits - 1
		letterIdxMax  = 63 / letterIdxBits
	)
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, l)
	for i, cache, remain := l-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(specifiedStr) {
			b[i] = specifiedStr[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b)) // nolint:gosec
}

// Int 指定范围内的随机数字
func Int(min int, max int) (int, error) {
	i, err := Int64(int64(min), int64(max))
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// Int32 指定范围内的随机数字
func Int32(min int32, max int32) (int32, error) {
	i, err := Int64(int64(min), int64(max))
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// Int64 指定范围内的随机数字
func Int64(min int64, max int64) (int64, error) {
	var (
		i   *big.Int
		err error
	)
	if max <= 0 {
		return 0, errors.New("max argument can not <= 0")
	}
	if max <= min {
		return 0, errors.New("max argument can not <= min argument")
	}
	if min == 0 {
		i, err = cRand.Int(cRand.Reader, big.NewInt(max))
		if err != nil {
			return 0, err
		}
		return i.Int64(), err
	}
	i, err = cRand.Int(cRand.Reader, big.NewInt(max-min))
	if err != nil {
		return 0, err
	}
	return min + i.Int64(), err
}
