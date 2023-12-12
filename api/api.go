package api

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid"

	"github.com/juancwu/hososuru/constants"
	"github.com/juancwu/hososuru/ws"
)

func CreateNewRoom(c echo.Context) error {
	// create new room id
	roomId := gonanoid.MustID(constants.RoomIdLen)

	// create new folder
	newPath := fmt.Sprintf("%s/%s", constants.TmpFileFolder, roomId)
	if err := os.MkdirAll(newPath, os.ModePerm); err != nil {
        log.Println(err)
		return err
	}

	file, err := c.FormFile("movie-upload")
	if err != nil {
        log.Println(err)
		return err
	}

	src, err := file.Open()
	if err != nil {
        log.Println(err)
		return err
	}
	defer src.Close()

    videoPath := fmt.Sprintf("%s/%s", newPath, file.Filename)
	dst, err := os.Create(videoPath)
	if err != nil {
        log.Println(err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
        log.Println(err)
		return err
	}

	ws.PendingRooms[roomId] = file.Filename

    c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/room/%s", roomId))

    return c.String(http.StatusOK, "Done!")
}

func ServeHoso(c echo.Context) error {
	roomId := c.Param("roomId")
	filename, ok := ws.PendingRooms[roomId]
	if !ok {
		return errors.New(fmt.Sprintf("No hoso for room with id: %s", roomId))
	}
	videoPath := fmt.Sprintf("%s/%s/%s", constants.TmpFileFolder, roomId, filename)
    http.ServeFile(c.Response().Writer, c.Request(), videoPath)
    return nil
}
