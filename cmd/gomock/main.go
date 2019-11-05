package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/huandu/xstrings"
)

var headerTpl = `// Code generated by gengo/cmd/gomock. DO NOT EDIT.

package mock

`

var (
	packagesBase  = flag.String("packagesBase", "", "package go path base")
	packagesDir   = flag.String("packagesDir", "", "package dir")
	targetDoc     = flag.String("targetDoc", "mock/doc.go", "target doc go file path")
	targetPrefix  = flag.String("targetPrefix", "", "target package prefix")
	targetPackage = flag.String("targetPackage", "mock", "target package")
)

func main() {
	flag.Parse()

	l := &lookup{}
	l.scan(*packagesBase, *packagesDir)

	buf := bytes.Buffer{}
	_, _ = buf.WriteString(headerTpl)
	for _, i := range l.interfaces() {
		targetFolder := strings.TrimPrefix(fmt.Sprintf("%s%s", *targetPrefix, strings.TrimPrefix(i.importPath, *packagesBase)), "/")
		if *targetPackage == "" {
			*targetPackage = strings.ToLower(i.packageName)
		}
		targetFile := fmt.Sprintf("%s.go", xstrings.ToSnakeCase(i.name))
		_, _ = buf.WriteString(fmt.Sprintf("//go:generate mockgen -destination %s/%s -package %s -self_package=%s %s %s\n", targetFolder, targetFile, *targetPackage, i.importPath, i.importPath, i.name))
	}

	if err := os.MkdirAll(filepath.Dir(*targetDoc), 0755); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(*targetDoc, buf.Bytes(), 0755); err != nil {
		panic(err)
	}
}

type lookup struct {
	buf []*ityp
	sync.Mutex
	sync.WaitGroup
}

func (ff *lookup) interfaces() []*ityp {
	ff.Lock()
	defer ff.Unlock()

	sort.Sort(sortedItyp(ff.buf))
	return ff.buf
}

func (ff *lookup) scan(goPackageBase, dir string) {
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
		ff.Add(1)
		go func() {
			defer ff.Done()
			ff.find(goPackageBase, path)
		}()
	}

	ff.Wait()
}

func (ff *lookup) append(interfaces []*ityp) {
	ff.Lock()
	defer ff.Unlock()

	ff.buf = append(ff.buf, interfaces...)
}

func (ff *lookup) find(base, dir string) {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, dir, func(info os.FileInfo) bool {
		return true
	}, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	interfaces := make([]*ityp, 0, 16)
	for p, pkg := range pkgs {
		for fileName, file := range pkg.Files {

			for _, obj := range file.Scope.Objects {
				typ, ok := obj.Decl.(*ast.TypeSpec)
				if !ok {
					continue
				}

				_, ok = typ.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}

				if strings.ToLower(string(obj.Name[0])) != string(obj.Name[0]) {
					interfaces = append(interfaces, &ityp{
						packageName: p,
						importPath:  filepath.Join(base, filepath.Dir(fileName)),
						name:        obj.Name,
					})
				}
			}
		}
	}

	ff.append(interfaces)
}

type ityp struct {
	packageName string
	importPath  string
	name        string
}

type sortedItyp []*ityp

func (a sortedItyp) Len() int           { return len(a) }
func (a sortedItyp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortedItyp) Less(i, j int) bool { return a[i].importPath+a[i].name < a[j].importPath+a[j].name }
