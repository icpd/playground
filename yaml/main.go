package main

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

func main() {
	yml := `
spec:
  template:
    spec:
      containers:
        - name: app
          num: 1
          image: nginx:1.25 # keep
`
	f, err := parser.ParseBytes([]byte(yml), 0)
	if err != nil {
		panic(err)
	}

	// YAMLPath：$.spec.template.spec.containers[0].image
	path, err := yaml.PathString("$.spec.template.spec.containers[0].image")
	if err != nil {
		panic(err)
	}

	// 构造新值节点
	newNode, _ := yaml.ValueToNode("nginx:1.27")

	// 替换
	if err := path.ReplaceWithNode(f, newNode); err != nil {
		panic(err)
	}

	// 编码输出
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	defer enc.Close()
	if err := enc.Encode(f.Docs[0].Body); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	path, err = yaml.PathString("$.spec.template.spec.containers[0].num")
	if err != nil {
		panic(err)
	}

	node, err := path.ReadNode(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(node.String())
}
