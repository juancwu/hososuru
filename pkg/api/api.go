package api

import (
	// "errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid"

	"github.com/juancwu/hososuru/pkg/constants"
	"github.com/juancwu/hososuru/pkg/ws"
)

func CreateNewRoom(c echo.Context) error {
	// create new room id
	roomId := gonanoid.MustID(constants.RoomIdLen)

	// create new folder
	newPath := fmt.Sprintf("%s/%s", constants.TmpFileFolder, roomId)
	if err := os.MkdirAll(newPath, os.ModePerm); err != nil {
		return err
	}

	file, err := c.FormFile("movie")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(fmt.Sprintf("%s/%s", newPath, file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	ws.PendingRooms[roomId] = file.Filename

    c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/room/%s", roomId))

    return c.String(http.StatusOK, "Done!")
}

func ServeHoso(c echo.Context) error {
	roomId := c.Param("roomId")
	// filename, ok := ws.PendingRooms[roomId]
	// if !ok {
	// 	return errors.New(fmt.Sprintf("No hoso for room with id: %s", roomId))
	// }
    filename := "Barbie.2023.1080p.WEBRip.x264.AAC5.1-[YTS.MX].mp4"

	videoPath := fmt.Sprintf("%s/%s/%s", constants.TmpFileFolder, roomId, filename)
    http.ServeFile(c.Response().Writer, c.Request(), videoPath)
    return nil
}
