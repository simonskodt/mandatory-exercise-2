package utils

import "sync"

// Counter is a thread safe simple counter.
type Counter struct {
	value int // value is the current value of the Counter.
	mu    sync.Mutex
}

// Increment increments the counter by 1.
func (c *Counter) Increment() {
	defer c.mu.Unlock()
	c.mu.Lock()
	c.value++
}

// Reset the counter.
func (c *Counter) Reset()  {
	defer c.mu.Unlock()
	c.mu.Lock()
	c.value = 0
}

// Value returns the Counter's value
func (c *Counter) Value() int {
	return c.value
}

// NewCounter creates and returns a new Counter with its value starting at 0.
func NewCounter() *Counter {
	return &Counter{value: 0}
}
