package httpserver

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port       int
	gin        *gin.Engine
	controller Controller
}

func NewServer(port int, controller Controller) Server {
	return Server{port: port, controller: controller}
}

func (s *Server) Listen() {
	s.gin = gin.Default()

	s.startRouter()

	err := s.gin.Run(fmt.Sprintf(":%d", s.port))

	if err != nil {
		log.Panic().Err(err).Msg("failed to start server")
	}

	log.Info().Msgf("server started on port %d", s.port)
}
