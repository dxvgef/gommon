package slice

// CompareStrings 比较两个strings slice，返回add新增元素和remove移除的元素
func CompareStrings(old, new []string) (add []string, remove []string) {
	for k := range new {
		if !InStr(old, new[k]) {
			add = append(add, new[k])
		}
	}
	for k := range old {
		if !InStr(new, old[k]) {
			remove = append(remove, old[k])
		}
	}

	return
}
