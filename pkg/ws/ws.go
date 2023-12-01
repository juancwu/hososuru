package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
    upgrader = websocket.Upgrader{}
)

func Handle (c echo.Context) error {
    ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
    if err != nil {
        return err
    }
    defer ws.Close()

    for {
        err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
        if err != nil {
            c.Logger().Error(err)
        }

        _, msg, err := ws.ReadMessage()
        if err != nil {
            c.Logger().Error(err)
        }
        fmt.Printf("%s\n", msg)
    }
}