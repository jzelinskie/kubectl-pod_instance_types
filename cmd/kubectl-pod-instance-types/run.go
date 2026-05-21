package main

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/util"
)

func run(ctx context.Context, factory util.Factory, printFlags *genericclioptions.PrintFlags, allNamespaces bool, streams genericclioptions.IOStreams) error {
	namespace, _, err := factory.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}
	if allNamespaces {
		namespace = ""
	}

	config, err := factory.ToRESTConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	outputFormat := ""
	if printFlags.OutputFormat != nil {
		outputFormat = *printFlags.OutputFormat
	}
	if outputFormat != "" {
		return runNonTable(ctx, clientset, namespace, printFlags, streams)
	}
	return runTable(ctx, clientset, namespace, streams)
}

func runTable(ctx context.Context, clientset kubernetes.Interface, namespace string, streams genericclioptions.IOStreams) error {
	table, err := fetchPodsAsTable(ctx, clientset, namespace)
	if err != nil {
		return err
	}

	nodeNames := extractNodeNames(table)
	instanceTypes, err := resolveInstanceTypes(ctx, clientset, nodeNames)
	if err != nil {
		return fmt.Errorf("resolving instance types: %w", err)
	}

	augmentTable(table, nodeNames, instanceTypes)

	p := printers.NewTablePrinter(printers.PrintOptions{})
	return p.PrintObj(table, streams.Out)
}

func fetchPodsAsTable(ctx context.Context, clientset kubernetes.Interface, namespace string) (*metav1.Table, error) {
	data, err := clientset.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Resource("pods").
		SetHeader("Accept", "application/json;as=Table;v=v1;g=meta.k8s.io,application/json").
		DoRaw(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching pods: %w", err)
	}

	var table metav1.Table
	if err := json.Unmarshal(data, &table); err != nil {
		return nil, fmt.Errorf("decoding table response: %w", err)
	}
	return &table, nil
}

type podNodeSpec struct {
	Spec struct {
		NodeName string `json:"nodeName"`
	} `json:"spec"`
}

func extractNodeNames(table *metav1.Table) []string {
	names := make([]string, len(table.Rows))
	for i, row := range table.Rows {
		if row.Object.Raw == nil {
			continue
		}
		var p podNodeSpec
		if err := json.Unmarshal(row.Object.Raw, &p); err == nil {
			names[i] = p.Spec.NodeName
		}
	}
	return names
}

func augmentTable(table *metav1.Table, nodeNames []string, instanceTypes map[string]string) {
	table.ColumnDefinitions = append(table.ColumnDefinitions, metav1.TableColumnDefinition{
		Name: "Instance-Type",
		Type: "string",
	})
	for i := range table.Rows {
		instanceType := "<none>"
		if i < len(nodeNames) {
			if t, ok := instanceTypes[nodeNames[i]]; ok && t != "" {
				instanceType = t
			}
		}
		table.Rows[i].Cells = append(table.Rows[i].Cells, instanceType)
	}
}

func runNonTable(ctx context.Context, clientset kubernetes.Interface, namespace string, printFlags *genericclioptions.PrintFlags, streams genericclioptions.IOStreams) error {
	return fmt.Errorf("output format %q is not yet supported", *printFlags.OutputFormat)
}
