package main

import (
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	if err := newRootCmd(streams).Execute(); err != nil {
		os.Exit(1)
	}
}
