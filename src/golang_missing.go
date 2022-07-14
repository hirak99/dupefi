package main

import (
	"reflect"
	"testing"
)

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func Map[T any, U any](data []T, f func(T) U) []U {
	mapped := make([]U, len(data))
	for i, e := range data {
		mapped[i] = f(e)
	}
	return mapped
}

func Filter[T any](list []T, cond func(int, T) bool) []T {
	var newList []T
	for i, t := range list {
		if cond(i, t) {
			newList = append(newList, t)
		}
	}
	return newList
}

func FilterChan[T any](c <-chan T, f func(T) bool) <-chan T {
	out := make(chan T)
	go func() {
		for e := range c {
			if f(e) {
				out <- e
			}
		}
		close(out)
	}()
	return out
}

// Generator to slice.
func ChanToSlice[T any](c <-chan T) []T {
	var result []T
	for v := range c {
		result = append(result, v)
	}
	return result
}

// Testing methods.

func AssertEqual[T comparable](t *testing.T, got, want T) {
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func AssertSliceEqual[T any](t *testing.T, got, want []T) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("got %v, want %v", got, want)
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
		t.Fatalf("counts mismatch - got %v, want %v", got, want)
	}
}
