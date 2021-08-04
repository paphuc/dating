package socket

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomSocket struct {
	ID         primitive.ObjectID `json:"id"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *MessageSocket
	Private    bool `json:"private"`
}

// NewRoom creates a new Room
func NewRoom(id primitive.ObjectID, private bool) *RoomSocket {
	return &RoomSocket{
		ID:         id,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *MessageSocket),
		Private:    private,
	}
}

// id for Room Big
func IdBigRoom() primitive.ObjectID {
	roomHex, _ := primitive.ObjectIDFromHex("000000000000000000000000")
	return roomHex
}

// RunRoom runs our room, accepting various requests
func (room *RoomSocket) RunRoomSocket() {
	for {
		select {

		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			room.broadcastToClientsInRoom(message)
		}

	}
}

func (room *RoomSocket) registerClientInRoom(client *Client) {
	room.clients[client] = true
}

func (room *RoomSocket) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *RoomSocket) broadcastToClientsInRoom(message *MessageSocket) {

	if room.ID == IdBigRoom() {
		for client := range room.clients {
			if client.UserID == message.ReceiverID {
				client.send <- message
			}
		}
		return
	}

	for client := range room.clients {
		client.send <- message
	}
}

func (room *RoomSocket) GetId() primitive.ObjectID {
	return room.ID
}
