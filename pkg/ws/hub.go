package ws

import (
	"fmt"
	"log"
	"os"

	"github.com/juancwu/hososuru/pkg/constants"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var h = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

// map[roomId]filename
var PendingRooms = make(map[string]string)

func (h *hub) run() {
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
			}
			h.rooms[s.room][s.conn] = true
		case s := <-h.unregister:
			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
                    fmt.Println("here")
					if len(connections) == 0 {
						delete(h.rooms, s.room)
                        filename, ok := PendingRooms[s.room]
                        if ok {
                            delete(PendingRooms, s.room)

                            // delete tmp folder
                            tmpPath := fmt.Sprintf("%s/%s/%s", constants.TmpFileFolder, s.room, filename)
                            err := os.RemoveAll(tmpPath)
                            if err != nil {
                                log.Printf("Error: %v\n", err)
                            } else {
                                log.Printf("Temporary Folder Delete: %s", tmpPath)
                            }
                        }

					}
				}
			}
		case m := <-h.broadcast:
			connections := h.rooms[m.room]
            fmt.Println(m.data)
			for c := range connections {
				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.room)
					}
				}
			}
		}
	}
}
