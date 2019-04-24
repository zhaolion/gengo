package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/zhaolion/gengo/cmd/autogen/marshal-gen/generators"
	"k8s.io/gengo/args"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
	arguments := args.Default()

	if err := arguments.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		spew.Dump(err)
		klog.Errorf("Error: %v", err)
		os.Exit(1)
	}
}
