package slice

import "testing"

func TestCompareStrings(t *testing.T) {
	old := []string{"111", "aaa", "bbb", "ccc"}
	new := []string{"aaa", "ccc", "ddd", "eee"}
	add, remove := CompareStrings(old, new)
	t.Log(add)
	t.Log(remove)
}
