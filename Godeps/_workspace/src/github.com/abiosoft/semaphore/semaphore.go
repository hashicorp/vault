package semaphore

import (
	"sync"
	"time"
)

type Semaphore struct {
	permits int
	avail   int
	channel chan int
	aMutex  *sync.Mutex
	rMutex  *sync.Mutex
}

func New(permits int) *Semaphore {
	if permits < 1 {
		panic("Invalid number of permits. Less than 1")
	}
	return &Semaphore{
		permits,
		permits,
		make(chan int, permits),
		&sync.Mutex{},
		&sync.Mutex{},
	}
}

//Acquire one permit, if its not available the goroutine will block till its available
func (s *Semaphore) Acquire() {
	s.aMutex.Lock()
	s.channel <- 1
	s.avail--
	s.aMutex.Unlock()
}

//Similar to Acquire() but for many permits
func (s *Semaphore) AcquireMany(n int) {
	if n > s.permits {
		panic("To many requested permits")
	}
	s.aMutex.Lock()
	s.avail -= n
	for ; n > 0; n-- {
		s.channel <- 1
	}
	s.avail += n
	s.aMutex.Unlock()
}

//Similar to AcquireMany() but cancels if duration elapse before getting the permits.
//Returns true if successful and false if timeout occurs.
func (s *Semaphore) AcquireWithin(n int, d time.Duration) bool {
	timeout := make(chan bool, 1)
	cancel := make(chan bool, 1)
	go func() {
		time.Sleep(d)
		timeout <- true
	}()
	go func() {
		s.AcquireMany(n)
		timeout <- false
		if <-cancel {
			s.ReleaseMany(n)
		}
	}()
	if <-timeout {
		cancel <- true
		return false
	}
	cancel <- false
	return true
}

//Release one permit
func (s *Semaphore) Release() {
	s.rMutex.Lock()
	<-s.channel
	s.avail++
	s.rMutex.Unlock()
}

//Release many permits
func (s *Semaphore) ReleaseMany(n int) {
	if n > s.permits {
		panic("Too many requested releases")
	}
	for ; n > 0; n-- {
		s.Release()
	}
}

//Number of available unacquired permits
func (s *Semaphore) AvailablePermits() int {
	if s.avail < 0 {
		return 0
	}
	return s.avail
}

//Acquire all available permits and return the number of permits acquired
func (s *Semaphore) DrainPermits() int {
	n := s.AvailablePermits()
	if n > 0 {
		s.AcquireMany(n)
	}
	return n
}
