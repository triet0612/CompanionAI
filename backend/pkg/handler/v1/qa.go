package v1

import (
	"CompanionBackend/pkg/middlewares"
	"CompanionBackend/pkg/model"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// @Summary      create QA
// @Tags         Story
// @Produce      json
// @Param 		 id				path		string	true	"id" 		example(51eecb74-bd12-40b4-bd3d-71eaa2a7d71b)
// @Param 		 question		formData 	string	true	"question" 	example(What is a dog?)
// @Param		 attachment		formData	file	false 	"file"
// @Failure		 200		{object}	model.QA
// @Failure		 400		{object}	model.APIError
// @Failure		 404		{object}	model.APIError
// @Router       /story/{id}	[post]
func (h *Handler) createQA(c echo.Context) error {
	question := c.FormValue("question")
	if question == "" {
		return c.JSON(http.StatusBadRequest, model.APIError{
			Err: "file_error",
			Msg: "empty question",
		})
	}
	var attachment []byte
	var extension string
	file, err := c.FormFile("attachment")
	if err == nil {
		filename := strings.Split(file.Filename, ".")
		extension = filename[len(filename)-1]
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.APIError{
				Err: "file_error",
				Msg: "cannot open file",
			})
		}
		defer src.Close()
		attachment, err = io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.APIError{
				Err: "file_error",
				Msg: "cannot read file",
			})
		}
	}
	user := c.Get("user").(*jwt.Token).Claims.(*middlewares.CustomJWTClaims)
	story := model.Story{
		StoryID: c.Param("id"),
		UserID:  user.UserID,
	}
	row := h.db.QueryRow(context.Background(),
		"SELECT Context FROM Story WHERE StoryID = $1 AND UserID = $2",
		story.StoryID, story.UserID,
	)
	if err := row.Scan(&story.StoryContext); err != nil {
		log.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusBadRequest, model.APIError{
				Err: "prompt_err",
				Msg: "invalid chat_id, owner_id",
			})
		}
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "prompt_error",
			Msg: "invalid id",
		})
	}
	ollamaRes, err := h.ollamaGenerate(question, extension, story.StoryContext, attachment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "prompt_error",
			Msg: err.Error(),
		})
	}
	ans := model.QA{
		QAID:      uuid.NewString(),
		Question:  question,
		Answer:    ollamaRes.Response,
		Extension: extension,
	}
	story.StoryContext = ollamaRes.Context
	if _, err := h.db.Exec(context.Background(),
		"INSERT INTO QA VALUES ($1, $2, $3, $4, $5, $6)",
		story.StoryID, ans.QAID, ans.Question, ans.Answer, ans.Extension, attachment,
	); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "prompt_error",
			Msg: "unexpected server error",
		})
	}
	if _, err := h.db.Exec(context.Background(),
		"UPDATE Story SET Context = $1 WHERE StoryID = $2",
		story.StoryContext, story.StoryID,
	); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "prompt_error",
			Msg: "unexpected server error",
		})
	}
	return c.JSON(http.StatusOK, ans)
}

func (h *Handler) ollamaGenerate(question string, extension string, chatContext []int, attachment []byte) (*model.OllamaResponse, error) {
	ollamaBody := map[string]any{
		"model":   h.config.Dynamic["text-text-model"],
		"prompt":  question,
		"stream":  false,
		"context": chatContext,
	}
	if slices.Contains([]string{"png", "jpg"}, extension) {
		ollamaBody["model"] = h.config.Dynamic["image-text-model"]
		ollamaBody["images"] = []string{base64.StdEncoding.EncodeToString(attachment)}
	}
	resBody, err := json.Marshal(ollamaBody)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(
		h.config.LLM_URL+"/api/generate", "application/json",
		bytes.NewBuffer(resBody),
	)
	if err != nil {
		return nil, err
	}
	ollamaRes := model.OllamaResponse{}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&ollamaRes); err != nil {
		return nil, err
	}
	return &ollamaRes, nil
}
