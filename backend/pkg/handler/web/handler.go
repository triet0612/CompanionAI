package web

import (
	"CompanionBackend/pkg/config"
	"CompanionBackend/pkg/db"
	"CompanionBackend/template"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	db     *db.DBHelper
	config *config.Config
}

func Init(db *db.DBHelper, config *config.Config) *Handler {
	return &Handler{
		db:     db,
		config: config,
	}
}

func (h *Handler) Mount(e *echo.Echo) {
	e.Renderer = template.Init()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "home.html", nil)
	})
	e.GET("/login", func(c echo.Context) error {
		return c.Render(200, "login.html", nil)
	})
}
