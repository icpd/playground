/*
https://colobu.com/2024/08/13/string-bytes-benchmark/
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

// region unsafe
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

// region new unsafe go1.20
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

// region unsafe pointer
func toPointerBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func toPointerString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// endregion

var s = "hello, world"
var bts = []byte("hello, world")

func BenchmarkStringToBytes(b *testing.B) {
	var fns = map[string]func(string) []byte{
		"强制转换": toRawBytes,
		"传统转换": toReflectBytes,
		"新型转换": toBytes,
		"指针转换": toPointerBytes,
	}
	for name, fn := range fns {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bts = fn(s)
			}
		})
	}
}

func BenchmarkBytesToString(b *testing.B) {
	var fns = map[string]func([]byte) string{
		"强制转换": toRawString,
		"传统转换": toReflectString,
		"新型转换": toString,
		"指针转换": toPointerString,
	}
	for name, fn := range fns {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s = fn(bts)
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
