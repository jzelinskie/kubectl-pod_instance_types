package main

import (
	"reflect"
	"testing"

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
