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
