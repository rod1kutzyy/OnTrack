package http

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rod1kutzyy/OnTrack/internal/delivery/http/handlers"
)

type Server struct {
	engine *gin.Engine
	port   string
}

func NewServer(port string, todoHandler *handlers.TodoHandler) *Server {
	router := gin.Default()

	RegisterTodoRoutes(router, todoHandler)

	return &Server{
		engine: router,
		port:   port,
	}
}

func (s *Server) Run() error {
	address := fmt.Sprintf(":%s", s.port)
	log.Printf("Server started on port %s\n", s.port)
	return s.engine.Run(address)
}
