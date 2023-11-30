package ws

import (
    "github.com/gorilla/websocket"
)

type Client struct {
    Conn *websocket.Conn
    Id string
    RoomId string
    // TODO: add manager here
    Manager interface{}
    MsgChan chan string
}
