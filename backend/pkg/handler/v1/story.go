package v1

import (
	"CompanionBackend/pkg/middlewares"
	"CompanionBackend/pkg/model"
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// @Summary      get story list
// @Tags         Story
// @Produce      json
// @accept 		 json
// @Failure		 200		{array}		model.Story
// @Failure		 400		{object}	model.APIError
// @Failure		 404		{object}	model.APIError
// @Router       /story		[get]
func (h *Handler) getStoryList(c echo.Context) error {
	claims := c.Get("user").(*jwt.Token).Claims.(*middlewares.CustomJWTClaims)
	ans := []model.Story{}
	row, err := h.db.Query(context.Background(),
		"SELECT StoryID, CreationDate FROM Story WHERE UserID = $1",
		claims.UserID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Msg: "",
			Err: "",
		})
	}
	for row.Next() {
		temp := model.Story{}
		if err := row.Scan(&temp.StoryID, &temp.CreationDate); err != nil {
			continue
		}
		ans = append(ans, temp)
	}
	slog.Info(c.Request().Header.Get("Content-Type"))
	return c.JSON(http.StatusOK, ans)
}

// @Summary      create new story
// @Tags         Story
// @Produce      json
// @accept 		 json
// @Failure		 200		{object}	model.Story
// @Failure		 400		{object}	model.APIError
// @Failure		 404		{object}	model.APIError
// @Router       /story	[post]
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
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "chat_error",
			Msg: "unexpected server error",
		})
	}
	return c.JSON(http.StatusOK, ans)
}

// @Summary      get story detail
// @Tags         Story
// @Produce      json
// @Param 		 story_id	path		string	true	"story id" example(51eecb74-bd12-40b4-bd3d-71eaa2a7d71b)
// @Failure		 200		{object}	model.Story
// @Failure		 400		{object}	model.APIError
// @Failure		 404		{object}	model.APIError
// @Router       /story/{story_id}	[get]
func (h *Handler) getStoryByStoryID(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middlewares.CustomJWTClaims)
	id := c.Param("story_id")
	if uuid.Validate(id) != nil {
		return c.JSON(http.StatusOK, model.APIError{
			Err: "pair_error",
			Msg: "invalid id",
		})
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
			return c.JSON(http.StatusInternalServerError, model.APIError{
				Err: "pair_error",
				Msg: "unexpected server error",
			})
		}
	}
	ans.StoryID = id
	return c.JSON(http.StatusOK, ans)
}

// @Summary      delete story by story_id
// @Tags         Story
// @Produce      json
// @Param 		 story_id	path		string	true	"story id" example(51eecb74-bd12-40b4-bd3d-71eaa2a7d71b)
// @Failure		 200		{object}	nil
// @Failure		 400		{object}	model.APIError
// @Failure		 404		{object}	model.APIError
// @Router       /story/{story_id}		[delete]
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
			c.JSON(http.StatusInternalServerError, model.APIError{
				Err: "delete_story_error",
				Msg: "unexpected server error",
			})
		}
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "delete_story_error",
			Msg: "unexpected server error",
		})
	}
	if _, err = tx.Exec(ctx, "DELETE FROM QA WHERE StoryID = $1", id); err != nil {
		slog.Error(err.Error())
		if err = tx.Rollback(ctx); err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusInternalServerError, model.APIError{
				Err: "delete_story_error",
				Msg: "unexpected server error",
			})
		}
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "delete_story_error",
			Msg: "unexpected server error",
		})
	}
	if err = tx.Commit(ctx); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "delete_story_error",
			Msg: "unexpected server error",
		})
	}
	return c.JSON(http.StatusOK, "")
}

// @Summary      get qa image
// @Tags         Story
// @Produce      json
// @Param 		 id					path		string		true	"qa_id"
// @Failure		 200				{object}	model.Story
// @Failure		 400				{object}	model.APIError
// @Failure		 404				{object}	model.APIError
// @Router       /qa/image/{id}		[get]
func (h *Handler) getImageByStoryID(c echo.Context) error {
	id := c.Param("id")
	if uuid.Validate(id) != nil {
		return c.JSON(http.StatusBadRequest, model.APIError{
			Err: "image_err",
			Msg: "invalid id",
		})
	}
	var attachment []byte
	row := h.db.QueryRow(context.Background(),
		"SELECT Attachment FROM QA WHERE QAID = $1", id,
	)
	if err := row.Scan(&attachment); err != nil {
		slog.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, model.APIError{
				Err: "image_err",
				Msg: "no image",
			})
		}
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "image_err",
			Msg: "server error",
		})
	}
	return c.Blob(http.StatusOK, "image/jpeg", attachment)
}
