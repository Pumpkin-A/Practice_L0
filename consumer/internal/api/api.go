package api

import (
	"fmt"
	"practiceL0_go_mod/internal/bank"

	"github.com/gin-gonic/gin"
)

type Server struct {
	TransactionManager *bank.TransactionManager
	Router             *gin.Engine
}

func New(tm *bank.TransactionManager) (*Server, error) {
	s := &Server{
		Router:             gin.New(),
		TransactionManager: tm,
	}
	return s, nil
}

func (s *Server) Run() error {
	s.registerHandlers()

	err := s.runHTTPServer("localhost", 9090)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) registerHandlers() {
	r := s.Router.Group("/api")
	r.GET("/getOrder", s.HandleGetOrder)

}

func (s *Server) runHTTPServer(host string, port int) error {
	listenAddress := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("starting http listener at http://%s\n", listenAddress)

	return s.Router.Run(listenAddress)
}
