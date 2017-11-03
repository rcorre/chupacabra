package main

import "k8s.io/client-go/rest"
import "testing"

func TestKubeResources(t *testing.T) {
	k := kube{
		clients: map[string]rest.Interface{
			"foo": nil,
			"bar": nil,
			"baz": nil,
		},
	}

	expected := []string{
		"foo",
		"bar",
		"baz",
	}
	actual := k.Resources()

	if len(expected) != len(actual) {
		t.Errorf("Mismatched length: expected %v, got %v", expected, actual)
	}

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("Mismatched contents: expected %v, got %v", expected, actual)
		}
	}
}
