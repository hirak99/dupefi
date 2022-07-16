package sanity

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

func (s *setInternalStruct[T]) Count() int {
	return len(s.m)
}
