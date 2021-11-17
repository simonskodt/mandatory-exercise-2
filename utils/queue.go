package utils

import (
	"container/list"
	"sync"
)

// A Queue is a simple thread safe FIFO queue based on a doubly linked list.
type Queue struct {
	list *list.List // list is the doubly linked list containing the elements of the Queue.
	mu   sync.Mutex
}

// A tuple is the element of a Queue.
type tuple struct {
	lamport int32
	name    string
}

// Enqueue creates and adds a tuple to the back of the Queue.
func (q *Queue) Enqueue(lamport int32, name string) {
	defer q.mu.Unlock()
	q.mu.Lock()
	q.list.PushBack(&tuple{lamport: lamport, name: name})
}

// Dequeue returns and removes the first element of the Queue.
func (q *Queue) Dequeue() (int32, string) {
	defer q.mu.Unlock()
	q.mu.Lock()
	element := q.list.Front()
	q.list.Remove(element)
	return element.Value.(*tuple).lamport, element.Value.(*tuple).name
}

// IsEmpty Is the Queue empty?
func (q *Queue) IsEmpty() bool {
	return q.list.Len() == 0
}

// NewQueue creates and returns a new empty Queue.
func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}
