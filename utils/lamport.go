package utils

import "sync"

// A Lamport logical clock, which can be locked/unlocked.
type Lamport struct {
	T int32
	Mu sync.Mutex
}

// Increment increments the lamport clock by 1.
func (l *Lamport) Increment()  {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	l.T++
}

// MaxAndIncrement sets the lamport clock to the maximum value of itself and some other clock and increments it by 1
func (l *Lamport) MaxAndIncrement(other int32) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if l.T < other {
		l.T = other
	}

	l.T++
}
