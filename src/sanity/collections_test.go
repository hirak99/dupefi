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

func TestSet(t *testing.T) {
	s := MakeSet[string]()
	AssertEqual(t, s.Count(), 0)

	s.Add("hello")
	s.Add("world")

	AssertEqual(t, s.Count(), 2)
	Assert(t, s.Has("hello"))
	Assert(t, s.Has("world"))
	AssertEqual(t, s.HasInt("hello"), 1)
	Assert(t, !s.Has("foo"))
	AssertEqual(t, s.HasInt("foo"), 0)

	s.Remove("hello")
	s.Remove("unknown") // Doesn't exist - doesn't complain.

	AssertEqual(t, s.Count(), 1)
	Assert(t, !s.Has("hello"))
	Assert(t, s.Has("world"))
}
