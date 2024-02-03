package kocto

import "sync"

type ConcurrentBatcher[T any] struct {
	m     *sync.Cond
	data  []T
	size  int
	count int
}

func NewConcurrentBatcher[T any](size int) *ConcurrentBatcher[T] {
	return &ConcurrentBatcher[T]{
		m:     sync.NewCond(&sync.Mutex{}),
		data:  make([]T, size),
		size:  size,
		count: 0,
	}
}

// Add inserts an item into the batch, if the batch is full it will lock
func (b *ConcurrentBatcher[T]) Add(item T) bool {
	b.m.L.Lock()
	defer b.m.L.Unlock()

	for b.count == b.size {
		// Wait until there's space in the batch
		b.m.Wait()
	}

	b.data[b.count] = item
	b.count++

	return b.count == b.size
}

// Flush empties the batch, calling the provided function with the batched data.
// If there are any writes waithing, they will be allowed to complete after the batch is flushed.
func (b *ConcurrentBatcher[T]) Flush(f func([]T)) {
	b.m.L.Lock()

	f(b.data[:b.count])

	b.count = 0
	b.data = make([]T, b.size)

	b.m.Broadcast()
	b.m.L.Unlock()
}
