package api

import (
	"log"
	"net/http"
	"practiceL0_go_mod/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) HandleGetOrder(c *gin.Context) {
	var getOrderReq models.GetOrderReq
	orderUIDFromQuery := c.Query("OrderUID")
	if orderUIDFromQuery == "" {
		if err := c.ShouldBindJSON(&getOrderReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OrderUID is required"})
			return
		}
	} else {
		parsedUUID, err := uuid.Parse(orderUIDFromQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OrderUID format"})
			return
		}
		getOrderReq.UUID = parsedUUID
	}

	order, err := s.TransactionManager.GetOrderByUUID(getOrderReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(order)

	c.JSON(http.StatusOK, order)
}

func (s *Server) HandleGetOrderHTML(c *gin.Context) {
	// Получаем OrderUID из query параметра
	orderUID := c.DefaultQuery("OrderUID", "")
	if orderUID == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "OrderUID is required",
		})
		return
	}

	orderUUID, err := uuid.Parse(orderUID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "OrderUID is malformed",
		})
		return
	}

	order, err := s.TransactionManager.GetOrderByUUID(models.GetOrderReq{UUID: orderUUID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Отправляем данные в шаблон для рендера
	c.HTML(http.StatusOK, "order.html", gin.H{
		"Order": order,
	})
}
