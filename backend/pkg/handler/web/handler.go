package web

import (
	"CompanionBackend/pkg/config"
	"CompanionBackend/pkg/db"
	"CompanionBackend/pkg/middlewares"
	"CompanionBackend/template"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	e.POST("/login", h.login)
	e.GET("/login", func(c echo.Context) error {
		return c.Render(200, "login.html", nil)
	})
	e.POST("/register", h.register)
	e.GET("/register", func(c echo.Context) error {
		return c.Render(http.StatusOK, "register.html", nil)
	})

	e.GET("/story", h.getStoryList)
	e.POST("/story", h.createStory)
	e.GET("/story/:story_id", h.getStoryByStoryID)
	e.POST("/story/:story_id", h.createStoryQA)
	e.DELETE("/story/:story_id", h.deleteStoryByStoryID)

}

func (h *Handler) login(c echo.Context) error {
	type body struct {
		Email    string `form:"email"  example:"abc@gmail.com"`
		Password string `form:"password"  example:"refo"`
	}
	var u body
	if err := c.Bind(&u); err != nil {
		return c.HTML(http.StatusBadRequest, "field errors")
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return c.HTML(http.StatusBadRequest, "invalid email")
	}
	row := h.db.QueryRow(context.Background(),
		"SELECT UserID, (Password=crypt($1, Password)) AS Matched FROM UserAccount WHERE Email = $2",
		u.Password, u.Email,
	)
	var userID string
	var correct bool
	if err := row.Scan(&userID, &correct); err != nil {
		c.SetCookie(&http.Cookie{
			Name:     "jwt",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		})
		if errors.Is(err, pgx.ErrNoRows) {
			return c.HTML(http.StatusUnauthorized, "wrong email")
		}
		return c.HTML(http.StatusUnauthorized, "unexpected server error")
	}
	if !correct {
		c.SetCookie(&http.Cookie{
			Name:     "jwt",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		})
		return c.HTML(http.StatusUnauthorized, "wrong email or password")
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod(h.config.JWT_AUTH_METHOD), middlewares.CustomJWTClaims{
		UserID: userID,
		Email:  u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	})
	tokenString, err := token.SignedString(h.config.JWT_SECRET)
	if err != nil {
		return c.HTML(http.StatusUnauthorized, err.Error())
	}
	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(1000 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	})
	c.Response().Header().Set("HX-redirect", "/")
	return c.HTML(http.StatusOK, "Ok")
}

func (h *Handler) register(c echo.Context) error {
	type body struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}
	var u body
	if err := c.Bind(&u); err != nil {
		return c.HTML(http.StatusBadRequest, "field invalid")
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return c.HTML(http.StatusBadRequest, "invalid email")
	}
	id := uuid.New()
	if _, err := h.db.Exec(context.Background(),
		`INSERT INTO UserAccount VALUES ($1, $2, crypt($3, gen_salt('bf')), current_timestamp);`,
		id, u.Email, u.Password,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return c.HTML(http.StatusBadRequest, "used email")
			}
		}
		slog.Error(err.Error())
		return c.HTML(http.StatusInternalServerError, "unexpected server error")
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod(h.config.JWT_AUTH_METHOD), middlewares.CustomJWTClaims{
		UserID: id.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	})
	tokenString, err := token.SignedString(h.config.JWT_SECRET)
	if err != nil {
		slog.Error(err.Error())
		return c.HTML(http.StatusInternalServerError, "unexpected server error")
	}
	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(1000 * time.Hour),
		HttpOnly: true,
	})
	return c.HTML(http.StatusOK, "success")
}
