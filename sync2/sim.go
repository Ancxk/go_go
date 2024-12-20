package sync2

import "time"

// Semaphore is a counting semaphore with the option to
// specify a timeout.
type Semaphore struct {
	slots   chan struct{}
	timeout time.Duration
}

// NewSemaphore creates a Semaphore. The count parameter must be a positive
// number. A timeout of zero means that there is no timeout.
func NewSemaphore(count int, timeout time.Duration) *Semaphore {
	sem := &Semaphore{
		slots:   make(chan struct{}, count),
		timeout: timeout,
	}
	for i := 0; i < count; i++ {
		sem.slots <- struct{}{}
	}
	return sem
}

// Acquire returns true on successful acquisition, and
// false on a timeout.
func (sem *Semaphore) Acquire() bool {
	if sem.timeout == 0 {
		<-sem.slots
		return true
	}
	tm := time.NewTimer(sem.timeout)
	defer tm.Stop()
	select {
	case <-sem.slots:
		return true
	case <-tm.C:
		return false
	}
}

// TryAcquire acquires a semaphore if it's immediately available.
// It returns false otherwise.
func (sem *Semaphore) TryAcquire() bool {
	select {
	case <-sem.slots:
		return true
	default:
		return false
	}
}

// Release releases the acquired semaphore. You must
// not release more than the number of semaphores you've
// acquired.
func (sem *Semaphore) Release() {
	sem.slots <- struct{}{}
}

// Size returns the current number of available slots.
func (sem *Semaphore) Size() int {
	return len(sem.slots)
}

