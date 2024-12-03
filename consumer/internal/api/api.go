package api

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/models"

	"github.com/gin-gonic/gin"
)

type Server struct {
	OrderManager OrderManager
	Router       *gin.Engine
	Srv          *http.Server
}

type OrderManager interface {
	GetOrderByUUID(req models.GetOrderReq) (*models.Order, error)
}

func New(cfg config.Config, om OrderManager) (*Server, error) {
	s := &Server{
		OrderManager: om,
		Router:       gin.Default(),
	}

	listenAddress := fmt.Sprintf("%s:%d", "0.0.0.0", cfg.Server.Port)
	s.Srv = &http.Server{
		Addr:    listenAddress,
		Handler: s.Router,
	}
	s.registerHandlers()
	return s, nil
}

func (s *Server) registerHandlers() {
	// Регистрируем функции для использования в шаблонах
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	// Загружаем шаблоны с подключением функций
	s.Router.SetFuncMap(funcMap)
	s.Router.LoadHTMLGlob("templates/*")

	r := s.Router.Group("/api")
	r.GET("/getOrder", s.HandleGetOrder)

	// Главная страница
	s.Router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Новый маршрут для рендера HTML
	s.Router.GET("/viewOrder", s.HandleGetOrderHTML)
}

func (s *Server) RunHTTPServer() error {
	slog.Info("starting http listener", "address", fmt.Sprintf("http://%s", s.Srv.Addr))
	return s.Srv.ListenAndServe()
}
