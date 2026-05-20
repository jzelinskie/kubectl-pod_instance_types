package main

import (
	"context"

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

	cmd := &cobra.Command{
		Use:          "kubectl-pod-instance-types [flags]",
		Short:        "List pods with their node's cloud instance type",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(context.Background(), factory, printFlags, allNamespaces, streams)
		},
	}

	configFlags.AddFlags(cmd.Flags())
	printFlags.AddFlags(cmd)
	cmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List pods in all namespaces")

	return cmd
}
