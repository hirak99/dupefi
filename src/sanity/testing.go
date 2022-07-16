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
