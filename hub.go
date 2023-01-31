package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Hub struct {
	sync.Mutex
	// In session Rooms.
	rooms map[*Room]bool
	// rooms   map[*Room]bool
	// Register requests from the clients.
	register chan *Client

	// room waiting for second player
	waitingRoom *Room
	// wwaiting bool
	waiting bool
	// Delete room when it is empty
	delete chan *Room
}

type gameMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	// data     int    `json:"data"`
	// moveList []int  `json:"moveList"`
}

func newHub() *Hub {
	return &Hub{
		register: make(chan *Client),
		rooms:    make(map[*Room]bool),
		delete:   make(chan *Room),
		waiting:  false,
	}
}

func (h *Hub) run() {
	for {
		select {

		case client := <-h.register:
			// if len(h.rooms) == 0 {
			// 	h.Lock()
			// 	defer h.Unlock()
			// 	room := newRoom(h, client)
			// 	client.room = room
			// 	room.players[client] = true
			// 	go room.run()
			// 	h.rooms[room] = true

			// 	mgs := &gameMessage{Event: "turn", Data: 1}
			// 	fmt.Println(mgs)
			// 	val, err := json.Marshal(mgs)
			// 	fmt.Println(val)
			// 	if err != nil {
			// 		fmt.Print(err.Error())
			// 	}
			// 	go client.writePump()
			// 	go client.readPump()
			// 	client.send <- val
			// 	h.waitingRoom = room
			// 	h.waiting = true
			// } else {

			if h.waiting {
				h.Lock()

				h.waiting = false
				room := h.waitingRoom
				h.Unlock()
				room.players[client] = true
				client.room = room

				go client.writePump()
				go client.readPump()
				i := 1
				for client := range room.players {

					mgs := &gameMessage{Event: "turn", Data: i}
					fmt.Println(mgs)
					val, err := json.Marshal(mgs)
					fmt.Println(val)
					if err != nil {
						fmt.Print(err.Error())
					}
					client.send <- val
					i++
				}

			} else {
				h.Lock()

				h.waiting = true
				room := newRoom(h, client)
				h.waitingRoom = room
				h.Unlock()
				client.room = room
				room.players[client] = true
				go room.run()
				h.rooms[room] = true
				go client.writePump()
				go client.readPump()

			}
			// }
		case room := <-h.delete:
			fmt.Print("hub delete request came")
			h.Lock()
			for p := range room.players {
				p.conn.Close()
			}
			delete(h.rooms, room)
			if len(h.rooms) == 0 {
				h.waiting = false
			}
			h.Unlock()
			// case client := <-h.unregister:

			// 	close(client.send)

			// 	// case message := <-h.broadcast:
			// 	// 	for client := range h.clients {
			// 	// 		select {
			// 	// 		case client.send <- message:
			// 	// 		default:
			// 	// 			close(client.send)
			// 	// 			delete(h.clients, client)
			// 	// 		}
			// 	// 	}
			// }
			fmt.Println("Total rooms ", len(h.rooms))
		}
	}

}
