package pages

import (
    "github.com/labstack/echo/v4"
	// "github.com/juancwu/hososuru/pkg/ws"
)

type RoomPageData struct {
    RoomId string
}

func Room(c echo.Context) error {
    roomId := c.Param("roomId")

    // verify room
    // if _, ok := ws.PendingRooms[roomId]; !ok {
    //     return c.Render(200, "not-found.html", nil)
    // }

    data := RoomPageData{
        RoomId: roomId,
    }

	return c.Render(200, "room.html", data)
}
