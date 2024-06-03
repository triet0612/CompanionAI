package v1

import (
	"CompanionBackend/pkg/middlewares"
	"CompanionBackend/pkg/model"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) createChat(c echo.Context) error {
	log.Println(c.Request().Cookies())
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.APIError{
			Err: "chat_error",
			Msg: "no user",
		})
	}
	claims := user.Claims.(*middlewares.CustomJWTClaims)

	ans := model.Chat{
		ChatID:       uuid.New(),
		OwnerID:      uuid.MustParse(claims.UserID),
		CreationDate: time.Now(),
	}
	if _, err := h.db.Exec(context.Background(),
		"INSERT INTO Chat VALUES ($1, $2, $3)",
		ans.ChatID, ans.OwnerID, ans.CreationDate,
	); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "chat_error",
			Msg: "unexpected server error",
		})
	}
	return c.JSON(http.StatusOK, ans)
}

func (h *Handler) chatRequest(c echo.Context) error {

	return nil
}
