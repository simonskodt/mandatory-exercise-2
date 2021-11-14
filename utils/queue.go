package utils

import "container/list"

type Tuple struct {
	Lamport int32
	Name string
}

type Queue struct {
	queue *list.List
}

func (q *Queue) Enqueue(lamport int32, name string) {
	q.queue.PushBack(Tuple{Lamport: lamport, Name: name})
}

func (q *Queue) Dequeue() *Tuple {
	element := q.queue.Front()
	q.queue.Remove(element)
	return element.Value.(*Tuple)
}

func (q *Queue) IsEmpty() bool {
	return q.queue.Len() == 0
}

func NewQueue() *Queue {
	return &Queue{
		queue: list.New(),
	}
}