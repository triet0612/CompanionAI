package template

import (
	"embed"
	"html/template"
	"io"
	"log"

	"github.com/labstack/echo/v4"
)

//go:embed *.html
var Template embed.FS

type TemplateRenderer struct {
	templates *template.Template
}

func Init() *TemplateRenderer {
	t, err := template.ParseFS(Template, "*.html")
	if err != nil {
		log.Fatal(err)
	}
	return &TemplateRenderer{
		templates: t,
	}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}
