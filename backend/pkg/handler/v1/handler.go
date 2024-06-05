package v1

import (
	"CompanionBackend/pkg/config"
	"CompanionBackend/pkg/db"
	"CompanionBackend/pkg/middlewares"
	"CompanionBackend/pkg/model"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"time"

	_ "CompanionBackend/docs"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
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

// @title		API Documentation
// @version     1.0
// @BasePath	/api/v1
func (h *Handler) Mount(g *echo.Group) {
	g.GET("/docs/*", echoSwagger.WrapHandler)

	g.POST("/login", h.login)
	g.POST("/register", h.register)

	g.GET("/story", h.getStoryList)
	g.POST("/story", h.createStory)
	g.GET("/story/:story_id", h.getStoryByStoryID)
	g.POST("/story/:story_id", h.createStoryQA)
	g.DELETE("/story/:story_id", h.deleteStoryByStoryID)
	g.GET("/qa/image/:id", h.getImageByStoryID)
}

// @Summary      get jwt, return in header and cookie
// @Tags         Authentication
// @Produce      json
// @Param 		 body		body 		v1.login.body	true	"body"
// @Failure		 200		{object}	nil
// @Failure		 400		{object}	model.APIError
// @Failure		 404		{object}	model.APIError
// @Router       /login		[post]
func (h *Handler) login(c echo.Context) error {
	type body struct {
		Email    string `json:"email" example:"abc@gmail.com"`
		Password string `json:"password" example:"refo"`
	}
	var u body
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, model.APIError{
			Err: "auth_error",
			Msg: "field errors",
		})
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return c.JSON(http.StatusBadRequest, model.APIError{
			Err: "auth_error",
			Msg: "invalid email",
		})
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
			return c.JSON(http.StatusUnauthorized, model.APIError{
				Err: "auth_error",
				Msg: "wrong email",
			})
		}
		return c.JSON(http.StatusUnauthorized, model.APIError{
			Err: "auth_error",
			Msg: "unexpected server error",
		})
	}
	if !correct {
		c.SetCookie(&http.Cookie{
			Name:     "jwt",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		})
		return c.JSON(http.StatusUnauthorized, model.APIError{
			Err: "auth_error",
			Msg: "wrong email or password",
		})
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
		return c.JSON(http.StatusUnauthorized, model.APIError{
			Err: "auth_error",
			Msg: err.Error(),
		})
	}
	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(1000 * time.Hour),
		HttpOnly: true,
	})
	return c.JSON(http.StatusOK, nil)
}

// @Summary      create new account, login after success
// @Tags         Authentication
// @Produce      json
// @Param 		 body		body 		v1.register.body	true	"body"
// @Failure		 200		{object}	nil
// @Failure		 400		{object}	model.APIError
// @Failure		 404		{object}	model.APIError
// @Router       /register	[post]
func (h *Handler) register(c echo.Context) error {
	type body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var u body
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, model.APIError{
			Err: "register_error",
			Msg: "field invalid",
		})
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return c.JSON(http.StatusBadRequest, model.APIError{
			Err: "register_error",
			Msg: "invalid email",
		})
	}
	id := uuid.New()
	if _, err := h.db.Exec(context.Background(),
		`INSERT INTO UserAccount VALUES ($1, $2, crypt($3, gen_salt('bf')), current_timestamp);`,
		id, u.Email, u.Password,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return c.JSON(http.StatusBadRequest, model.APIError{
					Err: "register_error",
					Msg: "used email",
				})
			}
		}
		slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "register_error",
			Msg: "unexpected server error",
		})
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
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "login_error",
			Msg: "unexpected server error",
		})
	}
	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(1000 * time.Hour),
		HttpOnly: true,
	})
	return c.JSON(http.StatusOK, map[string]string{
		"message": "success",
	})
}
