package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	rooms       map[*Room]bool
	register    chan *websocket.Conn
	waitingRoom struct {
		sync.Mutex
		room    *Room
		waiting bool
	}
	delete chan *Room
}

type gameMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func newHub() *Hub {
	return &Hub{
		register: make(chan *websocket.Conn, 100),
		rooms:    make(map[*Room]bool),
		delete:   make(chan *Room, 20),
	}
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			client := &Client{hub: h, conn: conn, send: make(chan []byte, 256)}
			h.waitingRoom.Lock()
			if h.waitingRoom.waiting {
				h.waitingRoom.waiting = false
				room := h.waitingRoom.room
				room.players[client] = true
				client.room = room
				go client.writePump()
				go client.readPump()
				i := 1
				for client := range room.players {
					msg := &gameMessage{Event: "turn", Data: i}
					val, err := json.Marshal(msg)
					if err != nil {
						// TODO could improve further if a player does not get there turn they could hold the game hostage
						fmt.Print(err.Error())
					} else {
						client.send <- val
					}
					i++
				}
			} else {
				h.waitingRoom.waiting = true
				room := newRoom(h)
				h.waitingRoom.room = room
				client.room = room
				room.players[client] = true
				go room.run()
				h.rooms[room] = true
				go client.writePump()
				go client.readPump()
			}
			h.waitingRoom.Unlock()
		case room := <-h.delete:
			h.waitingRoom.Lock()
			for p := range room.players {
				p.conn.Close()
			}
			delete(h.rooms, room)
			h.waitingRoom.Unlock()
		}
	}
}
