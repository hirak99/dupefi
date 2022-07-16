package sanity

import (
	"testing"
)

func expectError(t *testing.T, f func()) {
	panicAsError = true
	defer func() { _ = recover() }()
	defer func() { panicAsError = false }()
	f()
	printStack()
	t.Errorf("Shouldn't reach here")
}

func TestAssertEqual(t *testing.T) {
	AssertEqual(t, 1, 1)
	expectError(t, func() { AssertEqual(t, 1, 2) })
}

func TestAssertSliceEqual(t *testing.T) {
	AssertSliceEqual(t,
		[]string{},
		[]string{})
	AssertSliceEqual(t,
		[]string{"a", "bb", "ccc"},
		[]string{"a", "bb", "ccc"})

	expectError(t, func() {
		AssertSliceEqual(t,
			[]string{"a", "bb", "ccc"},
			[]string{})
	})
	expectError(t, func() {
		AssertSliceEqual(t,
			[]string{"a", "bb", "ccc"},
			[]string{"a", "bb", "ccc", "dddd"})
	})
	expectError(t, func() {
		AssertSliceEqual(t,
			[]string{"a", "bb", "ccc"},
			[]string{"", "a", "bb", "ccc"})
	})
	expectError(t, func() {
		AssertSliceEqual(t,
			[]string{"a", "bb", "ccc"},
			[]string{"bb", "a", "ccc"})
	})
}

func TestAssertSliceEqualUnordered(t *testing.T) {
	AssertSliceEqualUnordered(t,
		[]string{},
		[]string{})
	AssertSliceEqualUnordered(t,
		[]string{"a", "bb", "ccc"},
		[]string{"a", "bb", "ccc"})
	AssertSliceEqualUnordered(t,
		[]string{"a", "bb", "ccc"},
		[]string{"bb", "a", "ccc"})

	expectError(t, func() {
		AssertSliceEqualUnordered(t,
			[]string{"a", "bb", "ccc"},
			[]string{})
	})
	expectError(t, func() {
		AssertSliceEqualUnordered(t,
			[]string{"a", "bb", "ccc"},
			[]string{"a", "bb", "ccc", "dddd"})
	})
	expectError(t, func() {
		AssertSliceEqualUnordered(t,
			[]string{"a", "bb", "ccc"},
			[]string{"", "a", "bb", "ccc"})
	})
}
