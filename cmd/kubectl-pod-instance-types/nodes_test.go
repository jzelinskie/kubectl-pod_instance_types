package main

import (
	"context"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func makeNode(name string, labels map[string]string) *corev1.Node {
	return &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: name, Labels: labels}}
}

func TestResolveInstanceTypes(t *testing.T) {
	tests := []struct {
		name      string
		nodes     []*corev1.Node
		nodeNames []string
		want      map[string]string
		wantErr   bool
	}{
		{
			name:      "new label",
			nodes:     []*corev1.Node{makeNode("n1", map[string]string{labelInstanceType: "m5.large"})},
			nodeNames: []string{"n1"},
			want:      map[string]string{"n1": "m5.large"},
		},
		{
			name:      "beta label fallback",
			nodes:     []*corev1.Node{makeNode("n1", map[string]string{labelInstanceTypeBeta: "t3.small"})},
			nodeNames: []string{"n1"},
			want:      map[string]string{"n1": "t3.small"},
		},
		{
			name: "new label takes precedence over beta",
			nodes: []*corev1.Node{makeNode("n1", map[string]string{
				labelInstanceType:     "m5.large",
				labelInstanceTypeBeta: "t3.small",
			})},
			nodeNames: []string{"n1"},
			want:      map[string]string{"n1": "m5.large"},
		},
		{
			name:      "no instance type label",
			nodes:     []*corev1.Node{makeNode("n1", nil)},
			nodeNames: []string{"n1"},
			want:      map[string]string{"n1": ""},
		},
		{
			name:      "empty node name skipped",
			nodes:     []*corev1.Node{},
			nodeNames: []string{""},
			want:      map[string]string{},
		},
		{
			name:      "deduplicates node names",
			nodes:     []*corev1.Node{makeNode("n1", map[string]string{labelInstanceType: "m5.large"})},
			nodeNames: []string{"n1", "n1"},
			want:      map[string]string{"n1": "m5.large"},
		},
		{
			name:      "node get error bubbles up",
			nodes:     []*corev1.Node{},
			nodeNames: []string{"missing-node"},
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects := make([]runtime.Object, len(tt.nodes))
			for i, n := range tt.nodes {
				objects[i] = n
			}
			client := fake.NewSimpleClientset(objects...)
			got, err := resolveInstanceTypes(context.Background(), client, tt.nodeNames)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
