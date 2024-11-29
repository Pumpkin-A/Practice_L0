package api

import (
	"fmt"
	"html/template"
	"net/http"
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/models"

	"github.com/gin-gonic/gin"
)

type Server struct {
	TransactionManager TransactionManager
	Router             *gin.Engine
	Srv                *http.Server
}

type TransactionManager interface {
	GetOrderByUUID(req models.GetOrderReq) (*models.Order, error)
}

func New(cfg config.Config, tm TransactionManager) (*Server, error) {
	s := &Server{
		TransactionManager: tm,
		Router:             gin.Default(),
	}

	listenAddress := fmt.Sprintf("%s:%d", "localhost", cfg.Server.Port)
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
	fmt.Printf("starting http listener at http://%s\n", s.Srv.Addr)
	return s.Srv.ListenAndServe()
}
