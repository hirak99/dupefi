// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func TestAssert(t *testing.T) {
	Assert(t, true)
	expectError(t, func() { Assert(t, false) })
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
