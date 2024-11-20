package api

import (
	"log"
	"net/http"
	"practiceL0_go_mod/internal/models"

	"github.com/gin-gonic/gin"
)

func (s *Server) HandleGetOrder(c *gin.Context) {
	var getOrderReq models.GetOrderReq
	if err := c.ShouldBindJSON(&getOrderReq); err != nil {
		//TODO: добавить валидацию
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	order, err := s.TransactionManager.GetOrderByUUID(getOrderReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(order)

	c.JSON(http.StatusOK, order)

}
