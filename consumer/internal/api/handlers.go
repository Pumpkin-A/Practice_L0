package api

import (
	"encoding/json"
	"fmt"
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
			//TODO: добавить валидацию
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "OrderUID is required"})
				return
			}
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(order)

	c.JSON(http.StatusOK, order)
}

// TODO: подправить
func (s *Server) HandleGetOrderHTML(c *gin.Context) {
	// Получаем OrderUID из query параметра
	orderUID := c.DefaultQuery("OrderUID", "")
	if orderUID == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "OrderUID is required",
		})
		return
	}

	// Строим URL для API
	apiURL := fmt.Sprintf("http://localhost:9090/api/getOrder?OrderUID=%s", orderUID)

	// Отправляем запрос к API
	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": fmt.Sprintf("Failed to fetch order: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	// Парсим ответ в структуру
	var orderData models.Order
	if err := json.NewDecoder(resp.Body).Decode(&orderData); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": fmt.Sprintf("Error parsing response"),
		})
		return
	}

	// Отправляем данные в шаблон для рендера
	c.HTML(http.StatusOK, "order.html", gin.H{
		"Order": orderData,
	})
}
