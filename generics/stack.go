package generics

type Stack[T any] struct {
	nodes []*T
	Count int
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(n *T) {
	s.nodes = append(s.nodes[:s.Count], n)
	s.Count++
}

func (s *Stack[T]) Pop() *T {
	if s.Count == 0 {
		return nil
	}
	s.Count--
	return s.nodes[s.Count]
}
