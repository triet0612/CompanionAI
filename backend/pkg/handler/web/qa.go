package web

import (
	"CompanionBackend/pkg/middlewares"
	"CompanionBackend/pkg/model"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

func (h *Handler) createStoryQA(c echo.Context) error {
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
		StoryID: c.Param("story_id"),
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
	ans := model.QA{
		QAID:      uuid.NewString(),
		Question:  question,
		Answer:    "",
		Extension: extension,
	}
	if _, err := h.db.Exec(context.Background(),
		"INSERT INTO QA VALUES ($1, $2, $3, $4, $5, $6, current_timestamp)",
		story.StoryID, ans.QAID, ans.Question, ans.Answer, ans.Extension, attachment,
	); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "prompt_error",
			Msg: "unexpected server error",
		})
	}
	ollamaRes, err := h.ollamaGenerate(question, extension, story.StoryContext, attachment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Err: "prompt_error",
			Msg: err.Error(),
		})
	}
	story.StoryContext = ollamaRes.Context
	ans.Answer = ollamaRes.Response
	if _, err := h.db.Exec(context.Background(),
		"UPDATE QA SET Answer = $1 WHERE QAID = $2",
		ans.Answer, ans.QAID,
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
	img := ""
	if len(attachment) != 0 {
		img = fmt.Sprintf(`<img class="h-72 px-5 object-fill" src="/api/v1/qa/image/%s" alt="image">`, ans.QAID)
	}
	return c.HTML(http.StatusOK, fmt.Sprintf(`<li>
	<div class="chat chat-end">
		<div class="chat-bubble">
			%s
		</div>
	</div>
	<div class="chat chat-start">
		<div class="chat-bubble">
			%s
		</div>
	</div>
	%s
</li>`, ans.Question, ans.Answer, img))
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
