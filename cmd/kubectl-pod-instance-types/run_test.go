package main

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestAugmentTable(t *testing.T) {
	table := &metav1.Table{
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{Name: "Name"}, {Name: "Ready"}, {Name: "Status"}, {Name: "Restarts"}, {Name: "Age"},
		},
		Rows: []metav1.TableRow{
			{Cells: []interface{}{"pod-1", "1/1", "Running", 0, "1d"}},
			{Cells: []interface{}{"pod-2", "0/1", "Pending", 0, "5m"}},
		},
	}
	nodeNames := []string{"node-a", ""}
	instanceTypes := map[string]string{"node-a": "m5.large"}

	augmentTable(table, nodeNames, instanceTypes)

	if len(table.ColumnDefinitions) != 6 {
		t.Fatalf("expected 6 columns, got %d", len(table.ColumnDefinitions))
	}
	if table.ColumnDefinitions[5].Name != "Instance-Type" {
		t.Errorf("expected column %q, got %q", "Instance-Type", table.ColumnDefinitions[5].Name)
	}
	if table.Rows[0].Cells[5] != "m5.large" {
		t.Errorf("row 0: expected %q, got %v", "m5.large", table.Rows[0].Cells[5])
	}
	if table.Rows[1].Cells[5] != "<none>" {
		t.Errorf("row 1: expected %q for empty node name, got %v", "<none>", table.Rows[1].Cells[5])
	}
}

func TestExtractNodeNames(t *testing.T) {
	table := &metav1.Table{
		Rows: []metav1.TableRow{
			{Object: runtime.RawExtension{Raw: []byte(`{"spec":{"nodeName":"node-a"}}`)}},
			{Object: runtime.RawExtension{Raw: []byte(`{"spec":{}}`)}},
			{Object: runtime.RawExtension{Raw: nil}},
		},
	}

	got := extractNodeNames(table)
	want := []string{"node-a", "", ""}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAnnotatePods(t *testing.T) {
	pods := &corev1.PodList{
		Items: []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "pod-1"},
				Spec:       corev1.PodSpec{NodeName: "node-a"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "pod-2"},
				Spec:       corev1.PodSpec{NodeName: ""},
			},
		},
	}
	instanceTypes := map[string]string{"node-a": "m5.large"}
	annotatePods(pods, instanceTypes)

	if pods.Items[0].Annotations[labelInstanceType] != "m5.large" {
		t.Errorf("pod-1: expected annotation %q, got %q", "m5.large", pods.Items[0].Annotations[labelInstanceType])
	}
	if _, ok := pods.Items[1].Annotations[labelInstanceType]; ok {
		t.Errorf("pod-2: expected no annotation for pending pod, got %q", pods.Items[1].Annotations[labelInstanceType])
	}
}
