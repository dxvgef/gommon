package slice

// InStr 检查string值在一个string slice中是否存在
func InStr(s []string, str string) bool {
	for k := range s {
		if str == s[k] {
			return true
		}
	}
	return false
}

// InInt 检查int值在一个int slice中是否存在
func InInt(i []int, is int) bool {
	for k := range i {
		if is == i[k] {
			return true
		}
	}
	return false
}

// InInt64 检查int64值在一个int64 slice中是否存在
func InInt64(i []int64, is int64) bool {
	for k := range i {
		if is == i[k] {
			return true
		}
	}
	return false
}

// InFloat32 检查float32值在一个float32 slice中是否存在
func InFloat32(f []float32, fs float32) bool {
	for k := range f {
		if fs==f[k] {
			return true
		}
	}
	return false
}

// InFloat64 检查float64值在一个float64 slice中是否存在
func InFloat64(f []float64, fs float64) bool {
	for k := range f {
		if fs == f[k] {
			return true
		}
	}
	return false
}