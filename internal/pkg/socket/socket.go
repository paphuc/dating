package socket

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WsServer struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *MessageSocket
	Rooms      map[*RoomSocket]bool
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *MessageSocket),
		Rooms:      make(map[*RoomSocket]bool),
	}
}

func (server *WsServer) Run() {
	for {
		select {

		case client := <-server.Register:
			server.registerClient(client)

		case client := <-server.Unregister:
			server.unregisterClient(client)

		case message := <-server.Broadcast:
			server.broadcastToClients(message)
		}

	}
}

func (server *WsServer) registerClient(client *Client) {
	server.Clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.Clients[client]; ok {
		delete(server.Clients, client)
	}
}

func (server *WsServer) broadcastToClients(messageSocket *MessageSocket) {
	for client := range server.Clients {
		client.send <- messageSocket
	}
}

func (server *WsServer) findRoomByID(ID primitive.ObjectID) *RoomSocket {
	var foundRoom *RoomSocket

	for room := range server.Rooms {
		if room.GetId() == ID {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) createRoom(id primitive.ObjectID, private bool) *RoomSocket {
	room := NewRoom(id, private)
	go room.RunRoomSocket()
	server.Rooms[room] = true

	return room
}

func (server *WsServer) findClientByID(ID primitive.ObjectID) *Client {
	var foundClient *Client
	for client := range server.Clients {
		if client.ID == ID {
			foundClient = client
			break
		}
	}

	return foundClient
}
