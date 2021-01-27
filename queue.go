package go_gen_server

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	next *node
	val  interface{}
}

type queue struct {
	head, tail *node
}

func newQueue() *queue {
	q := &queue{}
	stub := &node{}
	q.head = stub
	q.tail = stub
	return q
}

// Push can be safely called from multiple goroutines
func (q *queue) Push(x interface{}) {
	n := new(node)
	n.val = x
	// current producer acquires head node
	prev := (*node)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(n)))

	// release node to consumer
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prev.next)), unsafe.Pointer(n))
}

// Pop must be called from a single, consumer goroutine
func (q *queue) Pop() interface{} {
	tail := q.tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next)))) // acquire
	if next != nil {
		q.tail = next
		v := next.val
		next.val = nil
		return v
	}
	return nil
}

// Empty must be called from a single, consumer goroutine
func (q *queue) Empty() bool {
	tail := q.tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next))))
	return next == nil
}
