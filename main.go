package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/juancwu/hososuru/api"
	"github.com/juancwu/hososuru/constants"
	"github.com/juancwu/hososuru/views"
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

	e.GET("/", func(c echo.Context) error {
		views.Landing().Render(context.Background(), c.Response().Writer)
		return nil
	})
	e.GET("/room/:roomId", func(c echo.Context) error {
		roomId := c.Param("roomId")

		filename, ok := ws.PendingRooms[roomId]
		if !ok {
			return c.Redirect(http.StatusTemporaryRedirect, "/404")
		}

		videoPath := fmt.Sprintf("%s/%s/%s", constants.TmpFileFolder, roomId, filename)
		mtype, err := mimetype.DetectFile(videoPath)
		if err != nil {
			return err
		}

		views.Room(roomId, mtype.String()).Render(context.Background(), c.Response().Writer)

		return nil
	})

	e.GET("/ws/:roomId", ws.Handle)

	e.POST("/api/new", api.CreateNewRoom)
	e.GET("/api/hoso/:roomId", api.ServeHoso)

	e.Any("/*", func(c echo.Context) error {
		views.NotFound().Render(context.Background(), c.Response().Writer)
		return nil
	})

	e.Logger.Fatal(e.Start(":5173"))
}
