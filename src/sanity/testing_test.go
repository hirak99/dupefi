package sanity

import (
	"testing"
)

func expectError(t *testing.T, f func()) {
	panicAsError = true
	defer func() { _ = recover() }()
	defer func() { panicAsError = false }()
	f()
	t.Errorf("Shouldn't reach here")
}

func TestAssertEqual(t *testing.T) {
	AssertEqual(t, 1, 1)
	expectError(t, func() { AssertEqual(t, 1, 2) })
}
