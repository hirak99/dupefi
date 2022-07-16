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
