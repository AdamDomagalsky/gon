package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sync/semaphore"
)
type Config struct {
	Port string
}

type Server struct {
	router   *http.ServeMux
	config   *Config
	counters *Counters
}

func NewServer(config *Config) *Server {
	Counter5 := Counter5(0)
	counters := &Counters{
		c1: &Counter1{value: 0},
		c2: &Counter2{value: 0},
		c3: &Counter3{value: atomic.Uint64{}},
		c4: &Counter4{
			Weighted: semaphore.NewWeighted(1),
			value:    0,
		},
		c5: Counter5,
		c6: &Counter6{
			value:           0,
			incrementAndGet: make(chan chan int64),
		},
	}
	go counters.c6.run()

	server := &Server{
		router:   http.NewServeMux(),
		config:   config,
		counters: counters,
	}
	server.setupRoutes()

	return server
}

func (s *Server) Start() error {
	server := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	log.Printf("Server started on port %s\n", s.config.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server gracefully stopped")
	return nil
}

func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}