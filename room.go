package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Room struct {
	sync.Mutex
	hub       *Hub
	players   map[*Client]bool
	broadcast chan gameMessage
	// Unregister requests from clients.
	unregister chan *Client
	waiting    bool
}

func newRoom(hub *Hub) *Room {
	return &Room{
		hub:        hub,
		players:    make(map[*Client]bool),
		waiting:    true,
		broadcast:  make(chan gameMessage),
		unregister: make(chan *Client),
	}
}

func (r *Room) run() {
	for {
		select {
		case message := <-r.broadcast:

			val, err := json.Marshal(message)
			fmt.Print("borad casting message ", message, "\n")
			if err != nil {
				fmt.Print(err.Error())
			}
			for client := range r.players {
				r.hub.Lock()
				select {
				case client.send <- val:
				default:
					close(client.send)
					delete(r.players, client)

				}
				r.hub.Unlock()
				if message.Event == "abandon" {
					r.hub.delete <- r
				}
			}

		case client := <-r.unregister:
			fmt.Println("unregister request came")
			client.conn.Close()
			delete(r.players, client)

		}
	}
}
