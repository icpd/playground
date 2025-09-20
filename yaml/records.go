package main

import (
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
	"github.com/k0kubun/pp/v3"
)

type PathVal struct {
	path string
	val  string
}

type Judge []*PathVal

var judge Judge = []*PathVal{
	{path: "$.records.cache.max_doc_size"},
	{path: "$.records.cache.ram_cache.size"},
	{path: "$.records.cache.ram_cache.algorithm"},
	{path: "$.records.cache.ram_cache_cutoff"},
}

func (j *Judge) needRestart(r io.Reader) bool {
	rst := false
	for _, pv := range *j {
		path, err := yaml.PathString(pv.path)
		if err != nil {
			panic(err)
		}

		node, err := path.ReadNode(r)
		if err != nil {
			fmt.Println("read node failed:", err, pv.path)
			continue
		}

		if node.String() != pv.val {
			pv.val = node.String()
			rst = true
		}
	}

	return rst
}

func main() {
	f, err := parser.ParseFile("records.yaml", 0)
	if err != nil {
		panic(err)
	}

	if judge.needRestart(f) {
		fmt.Println("need restart")
	}

	pp.Print(judge)
}
