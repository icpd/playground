/*
https://colobu.com/2024/08/13/string-bytes-benchmark/

goos: linux
goarch: amd64
cpu: AMD Ryzen 7 PRO 4750G with Radeon Graphics
BenchmarkStringToBytes/强制转换-16                 	22622428	        52.24 ns/op	      16 B/op	       1 allocs/op
BenchmarkStringToBytes/反射转换-16                 	469501591	         2.498 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringToBytes/新型转换-16                 	712188750	         1.583 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringToBytes/指针转换-16                 	714454291	         1.583 ns/op	       0 B/op	       0 allocs/op

goos: linux
goarch: amd64
cpu: AMD Ryzen 7 PRO 4750G with Radeon Graphics
BenchmarkBytesToString/强制转换-16                 	25351454	        43.93 ns/op	      16 B/op	       1 allocs/op
BenchmarkBytesToString/反射转换-16                 	673316130	         1.597 ns/op	       0 B/op	       0 allocs/op
BenchmarkBytesToString/新型转换-16                 	734822632	         1.533 ns/op	       0 B/op	       0 allocs/op
BenchmarkBytesToString/指针转换-16                 	704448451	         1.585 ns/op	       0 B/op	       0 allocs/op
*/
package main

import (
	"testing"
	"unsafe"
)

// region 强转
func toRawBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return []byte(s)
}

func toRawString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return string(b)
}

// endregion

// region 反射
func toReflectBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func toReflectString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// endregion

// region 新型 go1.20
func toBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
func toString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// endregion

// region 指针
func toPointerBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func toPointerString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// endregion

// benchmark
type tester[param, res any] struct {
	name string
	fn   func(param) res
}

var stringToBytesTests = []tester[string, []byte]{
	{"强制转换", toRawBytes},
	{"反射转换", toReflectBytes},
	{"新型转换", toBytes},
	{"指针转换", toPointerBytes},
}

var bytesToStringTest = []tester[[]byte, string]{
	{"强制转换", toRawString},
	{"反射转换", toReflectString},
	{"新型转换", toString},
	{"指针转换", toPointerString},
}

var s = "hello, world"
var bts = []byte("hello, world")

func BenchmarkStringToBytes(b *testing.B) {
	for _, t := range stringToBytesTests {
		b.Run(t.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bts = t.fn(s)
			}
		})
	}
}

func BenchmarkBytesToString(b *testing.B) {
	for _, t := range bytesToStringTest {
		b.Run(t.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s = t.fn(bts)
			}
		})
	}
}

// region go1.22 编译器强转类型优化
func BenchmarkStringToBytesRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = toRawBytes(s)
	}
}
func BenchmarkBytesToStringRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = toRawString(bts)
	}
}

// endregion
