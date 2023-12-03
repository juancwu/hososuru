package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/juancwu/hososuru/pages"
	"github.com/juancwu/hososuru/ws"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Failed to load env", err)
        os.Exit(1)
    }

    e := echo.New()
    e.Static("/static", "static")

    e.GET("/", pages.Index)
    e.GET("/ws", ws.Handle)

    e.Any("/*", func (c echo.Context) error {
        return c.Render(200, "not-found.html", nil)
    })

    e.Logger.Fatal(e.Start(":5173"))
}
