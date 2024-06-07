package web

import (
	"CompanionBackend/pkg/middlewares"
	"CompanionBackend/pkg/model"
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) getStoryList(c echo.Context) error {
	claims := c.Get("user").(*jwt.Token).Claims.(*middlewares.CustomJWTClaims)
	ans := []model.Story{}
	row, err := h.db.Query(context.Background(),
		"SELECT StoryID, CreationDate FROM Story WHERE UserID = $1",
		claims.UserID,
	)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "")
	}
	for row.Next() {
		temp := model.Story{}
		if err := row.Scan(&temp.StoryID, &temp.CreationDate); err != nil {
			continue
		}
		ans = append(ans, temp)
	}
	return c.Render(http.StatusOK, "story_list.html", ans)
}

func (h *Handler) createStory(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middlewares.CustomJWTClaims)

	ans := model.Story{
		StoryID:      uuid.NewString(),
		UserID:       claims.UserID,
		CreationDate: time.Now(),
	}
	if _, err := h.db.Exec(context.Background(),
		"INSERT INTO Story VALUES ($1, $2, $3, $4)",
		ans.StoryID, ans.UserID, ans.CreationDate, ans.StoryContext,
	); err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, "unexpected server error")
	}
	c.Response().Header().Set("HX-Redirect", "/story/"+ans.StoryID)
	return c.HTML(http.StatusOK, "")
}

func (h *Handler) getStoryByStoryID(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middlewares.CustomJWTClaims)
	id := c.Param("story_id")
	if uuid.Validate(id) != nil {
		return c.HTML(http.StatusBadRequest, "invalid id")
	}
	row := h.db.QueryRow(context.Background(), `SELECT s.StoryID, s.UserID, s.CreationDate, s.Context,
json_agg(JSON_BUILD_OBJECT(
	'qa_id', q.QAID,
	'question', q.Question,
	'answer', q.Answer,
	'extension', q.Extension,
	'creation_date', q.CreationDate
)) as QA
FROM QA q JOIN Story s ON q.StoryID = s.StoryID WHERE s.StoryID = $1 AND s.UserID = $2
GROUP BY s.StoryID`,
		id, claims.UserID,
	)
	ans := model.Story{}
	if err := row.Scan(&ans.StoryID, &ans.UserID, &ans.CreationDate, &ans.StoryContext, &ans.Content); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Println(err)
			return c.HTML(http.StatusInternalServerError, "unexpected server error")
		}
	}
	ans.StoryID = id
	sort.Slice(ans.Content, func(i, j int) bool {
		return ans.Content[i].CreationDate.Before(ans.Content[j].CreationDate)
	})
	return c.Render(http.StatusOK, "qa.html", ans)
}

func (h *Handler) deleteStoryByStoryID(c echo.Context) error {
	id := c.Param("story_id")
	if err := uuid.Validate(id); err != nil {
		return c.JSON(http.StatusBadRequest, model.APIError{
			Msg: "delete_story_error",
			Err: "invalid story_id",
		})
	}
	ctx := context.Background()
	claims := c.Get("user").(*jwt.Token).Claims.(*middlewares.CustomJWTClaims)
	tx, err := h.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "delete_story_error",
			Msg: "unexpected server error",
		})
	}
	if _, err = tx.Exec(ctx, "DELETE FROM Story WHERE StoryID = $1 AND UserID = $2",
		id, claims.UserID,
	); err != nil {
		slog.Error(err.Error())
		if err = tx.Rollback(ctx); err != nil {
			slog.Error(err.Error())
			c.HTML(http.StatusInternalServerError, "unexpected server error")
		}
		return c.HTML(http.StatusInternalServerError, "unexpected server error")
	}
	if _, err = tx.Exec(ctx, "DELETE FROM QA WHERE StoryID = $1", id); err != nil {
		slog.Error(err.Error())
		if err = tx.Rollback(ctx); err != nil {
			slog.Error(err.Error())
			c.HTML(http.StatusInternalServerError, "unexpected server error")
		}
		return c.HTML(http.StatusInternalServerError, "unexpected server error")
	}
	if err = tx.Commit(ctx); err != nil {
		slog.Error(err.Error())
		c.HTML(http.StatusInternalServerError, "unexpected server error")
	}
	c.Response().Header().Set("HX-Redirect", "/")
	return c.HTML(http.StatusOK, "")
}
