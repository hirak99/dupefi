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
	"fmt"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"
	"testing"
)

// Print stack trace, excluding this module.
func printStack() {
	// Unlike a sane language, debug.Stack() is a sequence of bytes.
	// So we resort to string matching to filter out this module.
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")
	var firstLineOut int
	for i, line := range lines {
		if line[0] != '\t' {
			continue
		}
		match, _ := regexp.MatchString(`/testing.go:`, line)
		// The first line after any `/testing.go` match is outside this module.
		if match {
			firstLineOut = i + 1
		} else if firstLineOut > 0 {
			break
		}
	}
	println("Partial stack trace -")
	println(strings.Join(lines[firstLineOut:], "\n"))
}

// For testing the tests.
var panicAsError bool

func logError(t *testing.T, msg string, a ...interface{}) {
	if panicAsError {
		panic(fmt.Sprintf(msg, a...))
	} else {
		printStack()
		t.Errorf(msg, a...)
	}
}

// Testing methods.

func Assert(t *testing.T, cond bool) {
	if !cond {
		logError(t, "condition failed")
	}
}

func AssertEqual[T comparable](t *testing.T, got, want T) {
	if got != want {
		logError(t, "got %v, want %v", got, want)
	}
}

func AssertSliceEqual[T any](t *testing.T, got, want []T) {
	if !reflect.DeepEqual(want, got) {
		logError(t, "got %v, want %v", got, want)
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
		logError(t, "counts mismatch - got %v, want %v", got, want)
	}
}
