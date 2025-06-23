package sets

import (
	"iter"
	"sync"
)

type Set[E comparable] struct {
	els map[E]struct{}
	mux sync.Mutex
}

func (s *Set[E]) Size() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.els == nil {
		s.els = map[E]struct{}{}
	}
	return len(s.els)
}

func (s *Set[E]) IsEmpty() bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.els == nil {
		s.els = map[E]struct{}{}
	}
	return len(s.els) == 0
}

func (s *Set[E]) Contains(el E) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.els == nil {
		s.els = map[E]struct{}{}
	}
	_, found := s.els[el]
	return found
}

func (s *Set[E]) valuesSeq(withLock bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		if withLock {
			s.mux.Lock()
			defer s.mux.Unlock()
		}
		for el := range s.els {
			if !yield(el) {
				return
			}
		}
	}
}

func (s *Set[E]) Values() iter.Seq[E] {
	return s.valuesSeq(true)
}

func (s *Set[E]) Add(els ...E) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.els == nil {
		s.els = map[E]struct{}{}
	}
	s.unsafeAdd(els...)
}

func (s *Set[E]) unsafeAdd(els ...E) {
	for _, el := range els {
		s.els[el] = struct{}{}
	}
}
