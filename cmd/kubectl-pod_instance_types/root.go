package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubectl/pkg/cmd/util"
)

func newRootCmd(streams genericclioptions.IOStreams) *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)
	factory := util.NewFactory(configFlags)
	printFlags := genericclioptions.NewPrintFlags("").WithTypeSetter(scheme.Scheme)

	var allNamespaces bool
	var selector string

	cmd := &cobra.Command{
		Use:          filepath.Base(os.Args[0]) + " [flags]",
		Short:        "List pods with their node's cloud instance type",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), factory, printFlags, allNamespaces, selector, streams)
		},
	}

	configFlags.AddFlags(cmd.Flags())
	printFlags.AddFlags(cmd)
	cmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List pods in all namespaces")
	cmd.Flags().StringVarP(&selector, "selector", "l", "", "Selector (label query) to filter on")

	return cmd
}
