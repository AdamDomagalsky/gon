package api

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/healthz", s.HealthCheckHandler)

	counter_path := "/api/v1/counter"
	s.router.HandleFunc(counter_path+"/bad", s.CounterBadAFHandler)
	s.router.HandleFunc(counter_path+"/mutex", s.CounterMutexHandler)
	s.router.HandleFunc(counter_path+"/atomic", s.CounterAtomicHandler)
	s.router.HandleFunc(counter_path+"/semaphore", s.CounterSemaphoreHandler)
	s.router.HandleFunc(counter_path+"/scoped", s.CounterScopedValueHandler)
	s.router.HandleFunc(counter_path+"/channel", s.CounterChannelHandler)
}
