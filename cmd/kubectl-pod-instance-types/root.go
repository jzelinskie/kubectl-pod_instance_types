package main

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func newRootCmd(streams genericclioptions.IOStreams) *cobra.Command {
	return &cobra.Command{Use: "kubectl-pod-instance-types"}
}
