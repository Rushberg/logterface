package utils

import (
	"errors"
	"unicode"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}

// Node represents a node in the linked list
type Node[T any] struct {
	Value T
	Next  *Node[T]
}

// Queue represents a generic queue implemented using a linked list
type Queue[T any] struct {
	Front *Node[T]
	Rear  *Node[T]
	Size  int // Track the size of the queue
}

// Enqueue adds an item to the end of the queue
func (q *Queue[T]) Enqueue(item T) {
	newNode := &Node[T]{Value: item, Next: nil}
	if q.Rear == nil {
		q.Front = newNode
		q.Rear = newNode
	} else {
		q.Rear.Next = newNode
		q.Rear = newNode
	}
	q.Size++ // Increment the size
}

// Dequeue removes and returns the item at the front of the queue
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T // Zero value for type T
	if q.Front == nil {
		return zero, errors.New("queue is empty")
	}
	item := q.Front.Value
	q.Front = q.Front.Next
	if q.Front == nil {
		q.Rear = nil
	}
	q.Size-- // Decrement the size
	return item, nil
}

// IsEmpty checks if the queue is empty
func (q *Queue[T]) IsEmpty() bool {
	return q.Front == nil
}

func (q *Queue[T]) ToSlice() []T {
	if q.Front == nil {
		return nil
	}

	// Preallocate slice to avoid repeated memory allocation
	res := make([]T, 0, q.Size)
	node := q.Front
	for node != nil {
		res = append(res, node.Value)
		node = node.Next
	}
	return res
}

// Length returns the number of items in the queue
func (q *Queue[T]) Length() int {
	return q.Size
}
