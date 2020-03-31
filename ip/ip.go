package ip

import (
	"bytes"
	"net"
)

/*
 检查IP地址是否在允许的范围
 startIP 起始IP
 endIP 截止IP
 objIP 被检查的IP
*/
func Between(startIP net.IP, endIP net.IP, objIP net.IP) bool {
	start16 := startIP.To16()
	end16 := endIP.To16()
	obj16 := objIP.To16()
	if start16 == nil || end16 == nil || obj16 == nil {
		return false
	}

	if bytes.Compare(obj16, start16) >= 0 && bytes.Compare(obj16, end16) <= 0 {
		return true
	}

	return false
}
