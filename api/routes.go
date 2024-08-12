package api

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/healthz", s.HealthCheckHandler)

	s.router.HandleFunc("/counter1", s.Counter1Handler)
	s.router.HandleFunc("/counter2", s.Counter2Handler)
	s.router.HandleFunc("/counter3", s.Counter3Handler)
	s.router.HandleFunc("/counter4", s.Counter4Handler)
	s.router.HandleFunc("/counter5", s.Counter5Handler)
	s.router.HandleFunc("/counter6", s.Counter6Handler)
}
