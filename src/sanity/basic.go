package sanity

import (
	"sort"
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

// Allows sorting by the values instead of the indices.
// Useful if you want to define the lessfn in a sane way, and use
// it for multiple sort calls.
func SaneSortSlice[T any](s []T, lessfn func(T, T) bool) {
	sort.Slice(s, func(i, j int) bool { return lessfn(s[i], s[j]) })
}
