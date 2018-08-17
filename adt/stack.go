package adt

//This stack is not thread-safe (all it needs is a mutex w/ locks and deferred unlocks).
type Stack struct {
	container []interface{}
}

func (s Stack) IsEmpty() bool {
	return len(s.container) == 0
}

func (s Stack) Size() int {
	return len(s.container)
}

func (s Stack) Peek() interface{} {
	if s.IsEmpty() {
		panic("Cannot peek from an empty stack!")
	}
	return s.container[len(s.container) - 1]
}

func (s *Stack) Push(item interface{}) {
	s.container = append(s.container, item)
}

func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		panic("Cannot pop from an empty stack!")
	}
	item := s.container[len(s.container) - 1]
	s.container = s.container[:(len(s.container) - 1)]
	return item
}