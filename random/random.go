package random

import (
	"bytes"
	"math/rand"
	"time"
	"unsafe"
)

// Upper 指定长度的随机大写字母
func Upper(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(Int32(65, 90)) != temp {
			temp = string(Int32(65, 90))
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
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(Int32(97, 122)) != temp {
			temp = string(Int32(97, 122))
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
	const letterIdxBits = 6
	const letterIdxMask = 1<<letterIdxBits - 1
	const letterIdxMax = 63 / letterIdxBits
	var src = rand.NewSource(time.Now().UnixNano())
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
func Int(min int, max int) int {
	rand.Seed(rand.Int63n(time.Now().UnixNano()))
	return min + rand.Intn(max-min)
}

// Int32 指定范围内的随机数字
func Int32(min int32, max int32) int32 {
	rand.Seed(rand.Int63n(time.Now().UnixNano()))
	return min + rand.Int31n(max-min)
}

// Int64 指定范围内的随机数字
func Int64(min int64, max int64) int64 {
	rand.Seed(rand.Int63n(time.Now().UnixNano()))
	return min + rand.Int63n(max-min)
}
