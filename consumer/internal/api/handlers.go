package api

import (
	"log/slog"
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
			slog.Error("GetOrder request", "method", "GET", "status", http.StatusBadRequest, "err", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "OrderUID is required"})
			return
		}
	} else {
		parsedUUID, err := uuid.Parse(orderUIDFromQuery)
		if err != nil {
			slog.Error("GetOrder request", "method", "GET", "status", http.StatusBadRequest, "err", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OrderUID format"})
			return
		}
		getOrderReq.UUID = parsedUUID
	}

	order, err := s.OrderManager.GetOrderByUUID(getOrderReq)
	if err != nil {
		slog.Error("GetOrder request", "method", "GET", "order", getOrderReq.UUID.String(), "status", http.StatusBadRequest, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slog.Info("GetOrder request", "method", "GET", "order", order.OrderUID, "status", http.StatusOK)

	c.JSON(http.StatusOK, order)
}

func (s *Server) HandleGetOrderHTML(c *gin.Context) {
	// Получаем OrderUID из query параметра
	orderUID := c.DefaultQuery("OrderUID", "")
	if orderUID == "" {
		slog.Error("GetOrder request", "method", "GET", "status", http.StatusBadRequest, "err", "OrderUID is required")
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "OrderUID is required",
		})
		return
	}

	orderUUID, err := uuid.Parse(orderUID)
	if err != nil {
		slog.Error("GetOrder request", "method", "GET", "status", http.StatusBadRequest, "err", "OrderUID is malformed")
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "OrderUID is malformed",
		})
		return
	}

	order, err := s.OrderManager.GetOrderByUUID(models.GetOrderReq{UUID: orderUUID})
	if err != nil {
		slog.Error("GetOrder request", "method", "GET", "status", http.StatusBadRequest, "order", orderUID, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slog.Info("GetOrder request", "method", "GET", "order", order.OrderUID, "status", http.StatusOK)
	// Отправляем данные в шаблон для рендера
	c.HTML(http.StatusOK, "order.html", gin.H{
		"Order": order,
	})
}
