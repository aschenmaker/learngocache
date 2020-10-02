package goCache

type ByteView struct {
	b []byte
}

// Len 返回byte的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回byte的数据拷贝切片slice
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String returns the data as a string type.
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
