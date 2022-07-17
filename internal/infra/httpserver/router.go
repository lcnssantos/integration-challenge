package httpserver

func (s *Server) startRouter() {
	s.gin.GET("/price/:currency", s.controller.Query)
}
