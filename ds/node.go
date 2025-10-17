package ds

type Node[T any] struct {
	// -1 means it has been visited
	ParentsNum int
	Children   []*Node[T]
	Val        *T
}

func (n Node[T]) travel(fn func(n Node[T])) {
	fn(n)
	for _, child := range n.Children {
		child.travel(fn)
	}
}
