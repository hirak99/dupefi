package main

import (
	"reflect"
	"sort"
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

// Allows sorting by the values instead of the indices.
// Useful if you want to define the lessfn in a sane way, and use
// it for multiple sort calls.
func SaneSortSlice[T any](s []T, lessfn func(*T, *T) bool) {
	sort.Slice(s, func(i, j int) bool { return lessfn(&s[i], &s[j]) })
}

// Set implementation.
type setInternalStruct[T comparable] struct {
	m map[T]bool
}

func MakeSet[T comparable]() setInternalStruct[T] {
	s := setInternalStruct[T]{}
	s.m = make(map[T]bool)
	return s
}

func (s *setInternalStruct[T]) Add(e T) {
	s.m[e] = true
}

func (s *setInternalStruct[T]) Has(e T) bool {
	_, ok := s.m[e]
	return ok
}

// Indicator function.
func (s *setInternalStruct[T]) HasInt(e T) int {
	_, ok := s.m[e]
	return If(ok, 1, 0)
}

func (s *setInternalStruct[T]) Remove(e T) {
	delete(s.m, e)
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
