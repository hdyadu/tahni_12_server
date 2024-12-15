package main

import (
	"encoding/json"
	"log"
	"sync"
)

type Room struct {
	sync.Mutex
	hub        *Hub
	players    map[*Client]bool
	broadcast  chan gameMessage
	unregister chan *Client
}

func newRoom(hub *Hub) *Room {
	return &Room{
		hub:        hub,
		players:    make(map[*Client]bool),
		broadcast:  make(chan gameMessage),
		unregister: make(chan *Client),
	}
}

func (r *Room) run() {
	for {
		select {
		case message := <-r.broadcast:
			val, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshaling message %v", err.Error())
			} else {
				log.Printf("game move %v", message.Data)
				for client := range r.players {
					select {
					case client.send <- val:
					default:
						close(client.send)
						delete(r.players, client)
					}
				}
			}

		case client := <-r.unregister:
			client.conn.Close()
			delete(r.players, client)
			if len(r.players) == 0 {
				r.hub.delete <- r
			}
		}
	}
}
