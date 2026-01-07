package realtime

import "encoding/json"

type Hub struct {
	// clients
	Clients map[*Client]bool
	// rooms
	// room "image123": client1, client2, client3
	// room "image777": client4, client5
	RoomsHub map[string]map[*Client]bool

	// register
	// hub will add it to clients list
	Register chan *Client

	// unregister
	Unregister chan *Client

	//broadcast channer
	Broadcast chan Event
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		RoomsHub:   make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Event),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.Clients[client] = true

		case client := <-h.Unregister:
			// clients se hata do
			delete(h.Clients, client)
			for room := range client.Rooms {
				// Jis bhi image id se ye client connected h vha se clent ko remove kar do
				delete(h.RoomsHub[room], client)
			}
			// close the channel for this client
			close(client.Send)

		case event := <-h.Broadcast:
			// Is event se related koi client hai
			if clients, ok := h.RoomsHub[event.RoomId]; ok {
				msg, _ := json.Marshal(event)
				for c := range clients {
					select {
					case c.Send <- msg:
					default:
						delete(h.Clients, c)
					}
				}
			}
		}
	}
}
