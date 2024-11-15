package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HandleGetTransaction(c *gin.Context) {
	log.Println("Hello")
	response := "ok"
	c.JSON(http.StatusOK, response)
}
