package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

type Counters struct {
	badAF       *CounterBadAF
	scopedValue func() uint64

	mutex     *CounterMutex
	atmoic    *CounterAtomic
	semaphore *CounterSemaphore
	channel   *CounterChannel
}

func NewAllCounters() *Counters {
	CounterScopedValue := CounterScopedValue(0)
	counters := &Counters{
		badAF:       &CounterBadAF{value: 0},
		scopedValue: CounterScopedValue,

		mutex:  &CounterMutex{value: 0},
		atmoic: &CounterAtomic{value: atomic.Uint64{}},
		semaphore: &CounterSemaphore{
			Weighted: semaphore.NewWeighted(1),
			value:    0,
		},
		channel: &CounterChannel{
			value:           0,
			incrementAndGet: make(chan chan uint64),
		},
	}
	go counters.channel.run()

	return counters
}

type CounterBadAF struct {
	value uint64 `default:"0"`
}

func (c *CounterBadAF) Increment() uint64 {
	c.value++
	return c.value
}

func (c *CounterBadAF) Get() uint64 {
	return c.value
}

type CounterMutex struct {
	sync.Mutex
	value uint64 `default:"0"`
}

func (c *CounterMutex) Increment() uint64 {
	c.Lock()
	defer c.Unlock()
	c.value++
	return c.value
}

func (c *CounterMutex) Get() uint64 {
	c.Lock()
	defer c.Unlock()
	return c.value
}

type CounterAtomic struct {
	value atomic.Uint64
}

func (c *CounterAtomic) Increment() uint64 {
	return c.value.Add(1)
}

func (c *CounterAtomic) Get() uint64 {
	return c.value.Load()
}

type CounterSemaphore struct {
	*semaphore.Weighted
	value uint64
}

func (c *CounterSemaphore) Increment() uint64 {
	c.Acquire(context.Background(), 1)
	c.value++
	defer c.Release(1)
	return c.value
}

func (c *CounterSemaphore) Get() uint64 {
	c.Acquire(context.Background(), 1)
	defer c.Release(1)
	return c.value
}

func CounterScopedValue(initialValue uint64) func() uint64 {
	count := initialValue

	return func() uint64 {
		count++
		return count
	}
}

type CounterChannel struct {
	incrementAndGet chan chan uint64
	value           uint64
}

func (c *CounterChannel) run() {
	for {
		select {
		case responseChan := <-c.incrementAndGet:
			c.value++
			responseChan <- c.value
		}
	}
}
func (c *CounterChannel) IncrementAndGet() uint64 {
	responseChan := make(chan uint64)
	c.incrementAndGet <- responseChan
	return <-responseChan
}

func (s *Server) CounterBadAFHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", s.counters.badAF.Increment())
}

func (s *Server) CounterMutexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", s.counters.mutex.Increment())
}

func (s *Server) CounterAtomicHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", s.counters.atmoic.Increment())
}

// func (s *Server) CounterSemaphoreHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	if err := s.counters.semaphore.Acquire(r.Context(), 1); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	s.counters.semaphore.value++
// 	fmt.Fprintf(w, "%d", s.counters.semaphore.value)
// 	s.counters.semaphore.Release(1)
// }

func (s *Server) CounterSemaphoreHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", s.counters.semaphore.Increment())
}

func (s *Server) CounterScopedValueHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", s.counters.scopedValue())
}

func (s *Server) CounterChannelHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", s.counters.channel.IncrementAndGet())
}
