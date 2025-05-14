package main

import (
	"fmt"
	"strings"

	"github.com/google/wire"
)

type Namer interface {
	Name() string
}

type Foo struct {
}

func (f Foo) Name() string {
	return "foo"
}

type Bar struct {
}

func (b Bar) Name() string {
	return "bar"
}

func NewFoo() Foo {
	return Foo{}
}

func NewBar() Bar {
	return Bar{}
}

func Names(names ...Namer) string {
	var b strings.Builder
	for i, n := range names {
		fmt.Println(n.Name())

		if i != 0 {
			b.WriteString(", ")
		}

		b.WriteString(n.Name())
	}

	return b.String()
}

// ProvideNamers 创建一个切片收集所有的Namer实现
func ProvideNamers(foo Foo, bar Bar) []Namer {
	return []Namer{foo, bar}
}

// ProvideNames 创建一个适配器函数，将切片转为可变参数调用
func ProvideNames(namers []Namer) string {
	return Names(namers...)
}

var provider = wire.NewSet(
	NewFoo,
	NewBar,
	ProvideNamers,
)

func main() {
	fmt.Println(app())
}
