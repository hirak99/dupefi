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

// Set implementation.
type set[T comparable] struct {
	m map[T]bool
}

func MakeSet[T comparable]() set[T] {
	s := set[T]{}
	s.m = make(map[T]bool)
	return s
}

func (s *set[T]) Add(e T) {
	s.m[e] = true
}

func (s *set[T]) Has(e T) bool {
	_, ok := s.m[e]
	return ok
}

// Indicator function.
func (s *set[T]) HasInt(e T) int {
	_, ok := s.m[e]
	return If(ok, 1, 0)
}

func (s *set[T]) Remove(e T) {
	delete(s.m, e)
}

func (s *set[T]) Count() int {
	return len(s.m)
}
