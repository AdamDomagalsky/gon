package api

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

type Counters struct {
	c1 *Counter1
	c2 *Counter2
	c3 *Counter3
	c4 *Counter4
	c5 func() int64
	c6 *Counter6
}

type Counter1 struct {
	value uint64 `default:"0"`
}

type Counter2 struct {
	sync.Mutex
	value uint64 `default:"0"`
}

type Counter3 struct {
	value atomic.Uint64
}

type Counter4 struct {
	*semaphore.Weighted
	value uint64
}

func Counter5(initialValue int64) func() int64 {
	count := initialValue

	return func() int64 {
		count++
		return count
	}
}

type Counter6 struct {
	incrementAndGet chan chan int64
	value           int64
}

func (c *Counter6) run() {
	for {
		select {
		case responseChan := <-c.incrementAndGet:
			c.value++
			responseChan <- c.value
		}
	}
}
func (c *Counter6) IncrementAndGet() int64 {
	responseChan := make(chan int64)
	c.incrementAndGet <- responseChan
	return <-responseChan
}

func (s *Server) Counter1Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	s.counters.c1.value++
	fmt.Fprintf(w, "%d", s.counters.c1.value)
}

func (s *Server) Counter2Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	s.counters.c2.Lock()
	defer s.counters.c2.Unlock()
	s.counters.c2.value++
	fmt.Fprintf(w, "%d", s.counters.c2.value)
}

func (s *Server) Counter3Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", s.counters.c3.value.Add(1))
}

func (s *Server) Counter4Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := s.counters.c4.Acquire(r.Context(), 1); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.counters.c4.value++
	fmt.Fprintf(w, "%d", s.counters.c4.value)
	s.counters.c4.Release(1)
}

func (s *Server) Counter5Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", s.counters.c5())
}

func (s *Server) Counter6Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", s.counters.c6.IncrementAndGet())
}
