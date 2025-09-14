package ws

type Hub struct {
	rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

func (h *Hub) GetRoom(roomID string) *Room {
	if r, ok := h.rooms[roomID]; ok {
		return r
	}
	room := NewRoom(roomID)
	h.rooms[roomID] = room
	go room.Run()
	return room
}
