package ws

import (
	"fmt"
	"time"
    "encoding/json"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
    upgrader = websocket.Upgrader{}
    PendingRooms = make(map[string]string)
)

type Message struct {
    Content string `json:"content"`
    EventType string `json:"eventType"`
}

type Client struct {
    Conn *websocket.Conn
    MsgChan chan string
}

func (c *Client) readMessage(ctx echo.Context) {
    err := c.Conn.SetReadDeadline(time.Now().Add(10 * time.Second))
    if err != nil {
        ctx.Logger().Error(err)
        return
    }

    c.Conn.SetPongHandler(func (appData string) error {
        err := c.Conn.SetReadDeadline(time.Now().Add(10 * time.Second))
        if err != nil {
            ctx.Logger().Error(err)
            return err
        }
        fmt.Println("pong")
        return nil
    })

    // defer c.Conn.Close()
    defer func() {
        fmt.Println("Closing ws connection...")
        c.Conn.Close()
    }()
    
    for {
        _, msg, err := c.Conn.ReadMessage()
        if err != nil {
            return
        }
        var message Message
        err = json.Unmarshal(msg, &message)
        if err != nil {
            ctx.Logger().Error(err)
            continue
        }
        fmt.Printf("From Client: %s\nEventType: %s\n", message.Content, message.EventType)
        c.MsgChan <- message.Content
    }
}

func (c *Client) writeMessage(ctx echo.Context) {
    defer c.Conn.Close()

    ticker := time.NewTicker(time.Second * 9)
    for {
        select {
        case text, ok := <-c.MsgChan:
            if !ok {
                return
            }
            fmt.Printf("Message to send: %s\n", text)
            err := c.Conn.WriteMessage(websocket.TextMessage, []byte(text))
            if err != nil{
                ctx.Logger().Error(err)
                return
            }
        case <- ticker.C:
            err := c.Conn.WriteMessage(websocket.PingMessage, []byte(""))
            if err != nil {
                ctx.Logger().Error(err)
                return
            }
            fmt.Println("Ping")
    }
    }
}

func Handle (c echo.Context) error {
    ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
    if err != nil {
        return err
    }

    client := Client{ws, make(chan string)}

    go client.readMessage(c)
    go client.writeMessage(c)

    return nil
}
