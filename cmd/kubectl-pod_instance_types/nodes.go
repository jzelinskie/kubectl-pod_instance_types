package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	labelInstanceType     = "node.kubernetes.io/instance-type"
	labelInstanceTypeBeta = "beta.kubernetes.io/instance-type"
)

func resolveInstanceTypes(ctx context.Context, client kubernetes.Interface, nodeNames []string) (map[string]string, error) {
	unique := make(map[string]struct{})
	for _, name := range nodeNames {
		if name != "" {
			unique[name] = struct{}{}
		}
	}

	result := make(map[string]string, len(unique))
	for name := range unique {
		node, err := client.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("getting node %q: %w", name, err)
		}
		result[name] = instanceTypeFromNode(node)
	}
	return result, nil
}

func instanceTypeFromNode(node *corev1.Node) string {
	if v, ok := node.Labels[labelInstanceType]; ok {
		return v
	}
	if v, ok := node.Labels[labelInstanceTypeBeta]; ok {
		return v
	}
	return ""
}
