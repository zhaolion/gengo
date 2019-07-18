package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

var headerTpl = `#!/bin/bash
set -euo pipefail
trap "echo 'error: Script failed: see failed command above'" ERR

`

var (
	packagesBase  = flag.String("packagesBase", "", "package go path base")
	packagesDir   = flag.String("packagesDir", "", "package dir")
	targetDoc     = flag.String("targetDoc", "gomock.sh", "target doc go file path")
	targetPrefix  = flag.String("targetPrefix", "mock", "target package prefix")
	targetPackage = flag.String("targetPackage", "mock", "target package")
)

func main() {
	flag.Parse()

	l := &lookup{}
	l.scan(*packagesBase, *packagesDir)

	buf := bytes.Buffer{}
	_, _ = buf.WriteString(headerTpl)
	for _, i := range l.interfaces() {
		targetFolder := strings.TrimPrefix(fmt.Sprintf("%s%s", *targetPrefix, strings.TrimPrefix(i.Name.Package, *packagesBase)), "/")
		if *targetPackage == "" {
			*targetPackage = strings.ToLower(i.Name.Name)
		}

		_, _ = buf.WriteString(fmt.Sprintf("mockgen -destination %s/%s -package %s -self_package=%s %s %s\n", targetFolder, "mock.go", *targetPackage, i.Name.Package, i.Name.Package, i.Name.Name))
	}

	if err := os.MkdirAll(filepath.Dir(*targetDoc), 0755); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(*targetDoc, buf.Bytes(), 0755); err != nil {
		panic(err)
	}
}

type lookup struct {
	sync.WaitGroup
	sync.Mutex
	buf []*types.Type
}

func (l *lookup) scan(goPackageBase, dir string) {
	paths := make(map[string]int)
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") {
			paths[filepath.Dir(path)] = 1
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	for path := range paths {
		path := path
		l.find(filepath.Join(goPackageBase, path))
	}
}

func (l *lookup) add(types []*types.Type) {
	l.Lock()
	defer l.Unlock()

	l.buf = append(l.buf, types...)
}

// find interface
func (l *lookup) find(dir string) {
	finder := parser.New()
	if err := finder.AddDir(dir); err != nil {
		panic(err)
	}

	universes, err := finder.FindTypes()
	if err != nil {
		panic(err)
	}

	result := make([]*types.Type, 0, 8)
	for p, pkg := range universes {
		if p != dir {
			continue
		}

		for _, t := range pkg.Types {
			if t.Kind == types.Interface {
				if strings.ToLower(string(t.Name.Name[0])) != string(t.Name.Name[0]) {
					result = append(result, t)
				}
			}
		}
	}

	l.add(result)
}

func (l *lookup) interfaces() []*types.Type {
	sort.Sort(sortedTypes(l.buf))
	return l.buf
}

type sortedTypes []*types.Type

func (a sortedTypes) Len() int           { return len(a) }
func (a sortedTypes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortedTypes) Less(i, j int) bool { return a[i].String() < a[j].String() }
