package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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
		page := views.Landing()
		views.RootLayout(page, nil).Render(context.Background(), c.Response().Writer)
		return nil
	})
	e.GET("/room/:roomId", func(c echo.Context) error {
		roomId := c.Param("roomId")

		filename, ok := ws.PendingRooms[roomId]
		if !ok {
			return errors.New(fmt.Sprintf("No hoso for room with id: %s", roomId))
		}

        videoPath := fmt.Sprintf("%s/%s/%s", constants.TmpFileFolder, roomId, filename)
        mtype, err := mimetype.DetectFile(videoPath)
        if err != nil {
            return err
        }

        room := views.Room(roomId, mtype.String())
        title := "Hososuru | Room " + roomId
        views.RootLayout(room, &title).Render(context.Background(), c.Response().Writer)

        return nil
	})
	e.GET("/ws", ws.Handle)

	e.POST("/api/new", api.CreateNewRoom)
	e.GET("/api/hoso/:roomId", api.ServeHoso)

	e.Any("/*", func(c echo.Context) error {
		component := views.NotFound()
		views.RootLayout(component, nil).Render(context.Background(), c.Response().Writer)
		return nil
	})

	e.Logger.Fatal(e.Start(":5173"))
}
