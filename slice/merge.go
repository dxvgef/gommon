package slice

import "bytes"

// Merge 多个[]byte合并成一个[]byte
func Merge(baseBytes []byte, bytesSlice ...[]byte) (result []byte, err error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	if _, err = buf.Write(baseBytes); err != nil {
		return
	}
	for k := range bytesSlice {
		if _, err = buf.Write(bytesSlice[k]); err != nil {
			return
		}
	}
	result = buf.Bytes()
	return
}
