package random

import (
	"bytes"
	"math/rand"
	"time"
)

// Upper 指定长度的随机大写字母
func Upper(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(Int32(65, 90)) != temp {
			temp = string(Int32(65, 90))
			result.WriteString(temp)
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
			result.WriteString(temp)
			i++
		}
	}
	return result.String()
}

// CustomString，指定长度的随机字符串，第二个参数限制只能出现指定的字符
func CustomString(l int, specifiedStr ...string) string {
	var tpl string
	if len(specifiedStr) > 0 {
		tpl = specifiedStr[0]
	} else {
		tpl = "abcdefghijklmnopqrstuwxyzABCDEFGHIJKLMNOPQRSTUWXYZ0123456789"
	}
	tplRunes := bytes.Runes([]byte(tpl))
	tplLen := len(tplRunes)
	resultRunes := make([]rune, l)
	for i := 0; i < l; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		resultRunes[i] = tplRunes[r.Intn(tplLen)]
	}
	return string(resultRunes)
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