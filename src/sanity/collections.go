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
