package sanity

import (
	"reflect"
	"testing"
)

// Testing methods.

func AssertEqual[T comparable](t *testing.T, got, want T) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func AssertSliceEqual[T any](t *testing.T, got, want []T) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func AssertSliceEqualUnordered[T comparable](t *testing.T, got, want []T) {
	getCounts := func(list []T) map[T]int {
		counts := make(map[T]int)
		for _, t := range list {
			if c, ok := counts[t]; ok {
				counts[t] = c + 1
			} else {
				counts[t] = 0
			}
		}
		return counts
	}
	if !reflect.DeepEqual(getCounts(want), getCounts(got)) {
		t.Errorf("counts mismatch - got %v, want %v", got, want)
	}
}
