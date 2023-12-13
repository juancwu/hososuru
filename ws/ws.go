package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/juancwu/hososuru/views"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid"
)

const (
	ChatEvent   = "chat_event"
	ToggleEvent = "toggle_event"
    PlaythroughEvent = "playthrough_event"
)

type messageEvent struct {
	Content string `json:"content"`
	EvtType string `json:"eventType"`
}

type broadcastEvent struct {
	clientId string
	msgEvt   messageEvent
}

type client struct {
	clientId string
	conn     *websocket.Conn
	msgChan  chan messageEvent
}

type room struct {
	clients    map[string]*client
	broadcast  chan broadcastEvent
	register   chan *client
	unregister chan string
}

var (
	upgrader     = websocket.Upgrader{}
	PendingRooms = make(map[string]string)
	rooms        = make(map[string]room)
)

func (c *client) readMessage(ctx echo.Context, roomId string) {
	err := c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		ctx.Logger().Error(err)
		return
	}

	_room, ok := rooms[roomId]
	if !ok {
		fmt.Printf("Read: ran read routine without proper room setup. Room ID: %s\n", roomId)
		c.conn.Close()
		return
	}

	c.conn.SetPongHandler(func(appData string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			ctx.Logger().Error(err)
			return err
		}
		fmt.Println("pong")
		return nil
	})

	// defer c.conn.Close()
	defer func() {
		fmt.Printf("Read: Closing connection to client (%s) in room (%s)\n", c.clientId, roomId)
		c.conn.Close()
		_room.unregister <- c.clientId
	}()

	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		var msgEvt messageEvent
		err = json.Unmarshal(rawMsg, &msgEvt)
		if err != nil {
			ctx.Logger().Error(err)
			continue
		}
		fmt.Printf("From client: %s\nEventType: %s\n", msgEvt.Content, msgEvt.EvtType)
		b := broadcastEvent{
			c.clientId,
			msgEvt,
		}
		_room.broadcast <- b
	}
}

func (c *client) writeMessage(ctx echo.Context, roomId string) {
	_room, ok := rooms[roomId]
	if !ok {
		fmt.Printf("Write: ran write routine without proper room setup. Room ID: %s\n", roomId)
		c.conn.Close()
		_room.unregister <- c.clientId
		return
	}

	defer func() {
		fmt.Printf("Write: Closing connection to client (%s) in room (%s)\n", c.clientId, roomId)
		c.conn.Close()
		_room.unregister <- c.clientId
	}()

	ticker := time.NewTicker(time.Second * 9)
	for {
		select {
		case msg, ok := <-c.msgChan:
			if !ok {
				fmt.Printf("Failed to retrieve data to write to client (%s)\n", c.clientId)
			} else {
				switch msg.EvtType {
				case ChatEvent:
					fmt.Printf("Received %s\n", ChatEvent)
                    message := views.Message(msg.Content)
                    buffer := views.ToBuffer(message)
                    err := c.conn.WriteMessage(websocket.TextMessage, buffer)
                    if err != nil {
                        ctx.Logger().Error(err)
                        return
                    }
                default:
					fmt.Printf("Received %s\n", msg.EvtType)
                    // json stringify
                    data, err := json.Marshal(msg)
                    if err != nil {
                        ctx.Logger().Error(err)
                    } else {
                        fmt.Printf("Message to send: %s\n", string(data))
                        err = c.conn.WriteMessage(websocket.TextMessage, data)
                        if err != nil {
                            ctx.Logger().Error(err)
                            return
                        }
                    }
				}
			}

		case <-ticker.C:
			err := c.conn.WriteMessage(websocket.PingMessage, []byte(""))
			if err != nil {
				ctx.Logger().Error(err)
				return
			}
			fmt.Println("Ping")
		}
	}
}

func (r room) run(ctx echo.Context, roomId string) {
	fmt.Printf("Room (%s) running\n", roomId)
	for {
		select {
		case b, ok := <-r.broadcast:
			if !ok {
				fmt.Printf("Broadcast to room (%s) failed\n", roomId)
			} else {
                if b.msgEvt.EvtType == PlaythroughEvent {
                    for clientId, c := range r.clients {
                        fmt.Printf("Client ID: %s\n", clientId)
                        if b.clientId != clientId {
                            c.msgChan <- b.msgEvt
                        }
                    }
                } else {
                    for _, c := range r.clients {
                        c.msgChan <- b.msgEvt
                    }
                }
			}

		case c, ok := <-r.register:
			if !ok {
				fmt.Printf("Register client to room (%s) failed\n", roomId)
			} else {
				r.clients[c.clientId] = c
			}

		case id, ok := <-r.unregister:
			if !ok {
				fmt.Printf("Unregister client to room (%s) failed\n", roomId)
			} else {
				_, ok = r.clients[id]
				if !ok {
					fmt.Printf("Client (%s) is not part of room (%s)\n", id, roomId)
				}

				delete(r.clients, id)

				n := len(r.clients)
				fmt.Printf("Remaining clients in room (%s): %d\n", roomId, n)
				if n == 0 {
					// remove room
					fmt.Printf("Removing room (%s) because its empty\n", roomId)
					delete(rooms, roomId)
					return
				}
			}
		}
	}
}

func Handle(c echo.Context) error {
	roomId := c.Param("roomId")
	_, ok := PendingRooms[roomId]
	if !ok {
		return errors.New("No room with given id")
	}

	_room, ok := rooms[roomId]
	if !ok {
		_room = room{
			make(map[string]*client),
			make(chan broadcastEvent),
			make(chan *client),
			make(chan string),
		}
		rooms[roomId] = _room
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	clientId := gonanoid.MustID(12)
	conn := client{clientId, ws, make(chan messageEvent)}

	go conn.readMessage(c, roomId)
	go conn.writeMessage(c, roomId)
	if !ok {
		// start room only when it is newly created
		go _room.run(c, roomId)
	}

	_room.register <- &conn

	return nil
}
