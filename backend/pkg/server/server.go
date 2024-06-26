package server

import (
	"CompanionBackend/assets"
	"CompanionBackend/pkg/config"
	"CompanionBackend/pkg/db"
	v1 "CompanionBackend/pkg/handler/v1"
	"CompanionBackend/pkg/handler/web"
	"CompanionBackend/pkg/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	db  *db.DBHelper
	cfg *config.Config
	e   *echo.Echo
}

func Init(db *db.DBHelper, cfg *config.Config) *Server {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${status} ${method}\turi=${uri}\t${error}\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8000", "http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	}))
	e.Use(middleware.Recover())
	m := middlewares.Middleware{
		JWT_SECRET:      cfg.JWT_SECRET,
		JWT_AUTH_METHOD: cfg.JWT_AUTH_METHOD,
	}

	e.Use(m.AuthMiddleware())
	e.StaticFS("/assets", assets.Assets)

	return &Server{
		cfg: cfg,
		e:   e,
		db:  db,
	}
}

func (s *Server) Run() {
	apiHandler := v1.Init(s.db, s.cfg)
	webHandler := web.Init(s.db, s.cfg)

	apiHandler.Mount(s.e)
	webHandler.Mount(s.e)

	s.e.Logger.Fatal(s.e.Start(":" + s.cfg.API_PORT))
}
