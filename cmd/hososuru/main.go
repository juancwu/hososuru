package main

import (
	"html/template"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/juancwu/hososuru/pkg/pages"
	"github.com/juancwu/hososuru/pkg/ws"
)


type TemplateRenderer struct {
    templates *template.Template
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Failed to load env", err)
        os.Exit(1)
    }

    templates, err := template.New("").ParseGlob("public/views/*.html")
    if err != nil {
        log.Fatalf("Error initializing templates: %v", err)
        os.Exit(1)
    }

    e := echo.New()
    e.Renderer = &TemplateRenderer{
        templates: templates,
    }

    e.Static("/static", "static")

    e.GET("/", pages.Index)
    e.GET("/ws", ws.Handle)

    e.Any("/*", func (c echo.Context) error {
        return c.Render(200, "not-found.html", nil)
    })

    e.Logger.Fatal(e.Start(":5173"))
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}
