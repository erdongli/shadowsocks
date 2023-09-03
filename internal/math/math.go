package math

import "cmp"

func Max[T cmp.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func Min[T cmp.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func IncrLittleEndian(b []byte) {
	for i := range b {
		b[i]++
		if b[i] == 1 {
			return
		}
	}
}
