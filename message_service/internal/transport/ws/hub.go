package ws

import "message_service/internal/service"

type Hub struct {
	rooms   map[string]*Room
	service service.MessageService
}

func NewHub(service service.MessageService) *Hub {
	return &Hub{
		rooms:   make(map[string]*Room),
		service: service,
	}
}

func (h *Hub) GetRoom(roomID string) *Room {
	if room, ok := h.rooms[roomID]; ok {
		return room
	}
	room := NewRoom(roomID, h.service)
	h.rooms[roomID] = room
	go room.Run()
	return room
}
