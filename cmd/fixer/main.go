package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(` // import "([\s\S]*?)"`)
var dir = flag.String("dir", "", "remove unused import comment")

func main() {
	flag.Parse()
	err := filepath.Walk(*dir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") {
			bs, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}

			if err := ioutil.WriteFile(path, reg.ReplaceAll(bs, []byte("")), 0755); err != nil {
				panic(err)
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
