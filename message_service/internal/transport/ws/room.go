package ws

import (
	"fmt"
	"message_service/internal/service"
)

type Room struct {
	ID         string
	clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	unregister chan *Client
	service    service.MessageService
}

func NewRoom(id string, service service.MessageService) *Room {
	return &Room{
		ID:         id,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 512),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		service:    service,
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.clients[client] = true
			fmt.Println("Client connected to room:", r.ID)

		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
				fmt.Println("Client disconnected from room:", r.ID)
			}

		case msg := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}

func (r *Room) BroadcastMessage(msg []byte) {
	select {
	case r.broadcast <- msg:
	default:
		fmt.Println("Broadcast channel full, message dropped for room:", r.ID)
	}
}
