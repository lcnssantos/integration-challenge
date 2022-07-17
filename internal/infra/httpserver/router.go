package httpserver

func (s *Server) startRouter() {
	s.gin.GET("/price/:currency", s.controller.Query)
	s.gin.POST("/service-c/callback", s.controller.Subscribe)
}
