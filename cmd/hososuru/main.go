package main

import (
	"html/template"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/juancwu/hososuru/pkg/pages"
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

    e.Use(middleware.Logger())
    e.Static("/static", "static")

    e.GET("/", pages.Index)

    e.Logger.Fatal(e.Start(":5173"))
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}
