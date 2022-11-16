package kocto

type Batcher[T any] struct {
	data  []T
	size  int
	count int
}

func NewBatcher[T any](size int) *Batcher[T] {
	return &Batcher[T]{
		data:  make([]T, size),
		size:  size,
		count: 0,
	}
}

func (b *Batcher[T]) Add(item T) bool {
	b.data[b.count] = item
	b.count++

	return b.count == b.size
}

func (b *Batcher[T]) Flush(f func([]T)) {
	f(b.data[:b.count])

	b.count = 0
	b.data = make([]T, b.size)
}
