package main

import (
	"unicode/utf8"
)

// https://github.com/golang/go/blob/9e9b1f57c26a6d13fdaebef67136718b8042cdba/src/cmd/compile/internal/walk/range.go#L292
func main() {
	s := "EN,中文"

	ha := s
	for hv1 := 0; hv1 < len(ha); {
		hv1t := hv1
		hv2 := rune(ha[hv1])
		if hv2 < utf8.RuneSelf {
			hv1++
		} else {
			var size int
			hv2, size = utf8.DecodeRune([]byte(ha[hv1:]))
			hv1 += size
		}

		println(hv1t, string(hv2))
	}

	println()
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		println(i, string(rs[i]))
	}
}
