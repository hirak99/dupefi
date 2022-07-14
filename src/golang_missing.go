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
