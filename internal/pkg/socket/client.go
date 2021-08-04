package socket

import (
	"dating/internal/app/api/types"
	"dating/internal/pkg/glog"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID       primitive.ObjectID
	RoomId   primitive.ObjectID
	UserID   primitive.ObjectID
	wsServer *WsServer
	conn     *websocket.Conn
	send     chan *MessageSocket
	save     *chan SaveMessage
	rooms    map[*RoomSocket]bool
}

func NewClient(conn *websocket.Conn, wsServer *WsServer, idRoom, idUser primitive.ObjectID, sm *chan SaveMessage) *Client {
	return &Client{
		ID:       primitive.NewObjectID(),
		RoomId:   idRoom,
		UserID:   idUser,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan *MessageSocket),
		rooms:    make(map[*RoomSocket]bool),
		save:     sm,
	}

}

func (c *Client) Read(logger glog.Logger) {
	defer c.conn.Close()
	for {
		var mgs *MessageSocket
		err := c.conn.ReadJSON(&mgs)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf("Failed when read message  %v", err)
			}
			return
		}

		if mgs.Attachments == nil {
			mgs.Attachments = []string{}
		}
		c.handleNewMessage(mgs)
	}
}

func (c *Client) Write(logger glog.Logger) {
	defer c.conn.Close()
	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			logger.Errorf("Failed when read message  %v", err)
			return
		}
	}
}

func (client *Client) handleNewMessage(jsonMessage *MessageSocket) {
	roomID := client.RoomId
	jsonMessage.RoomID = roomID

	switch jsonMessage.Action {
	case SendMessageAction:

		if room := client.wsServer.findRoomByID(roomID); room != nil {

			jsonMessage.ID = primitive.NewObjectID()

			room.broadcast <- jsonMessage

			if roomBig := client.wsServer.findRoomByID(IdBigRoom()); room != nil {
				roomBig.broadcast <- jsonMessage
			}

			sm := &SaveMessage{
				message: &types.Message{
					ID:          jsonMessage.ID,
					RoomID:      roomID,
					ReceiverID:  jsonMessage.ReceiverID,
					SenderID:    jsonMessage.SenderID,
					Content:     jsonMessage.Content,
					Attachments: jsonMessage.Attachments,
					CreateAt:    jsonMessage.CreateAt,
				},
			}

			if roomID != IdBigRoom() {
				*client.save <- *sm
			}
		}

	case JoinRoomAction:
		client.handleJoinRoomMessage(*jsonMessage)
	case LeaveRoomAction:
		client.handleLeaveRoomMessage(*jsonMessage)

	}

}

func (client *Client) handleLeaveRoomMessage(message MessageSocket) {
	room := client.wsServer.findRoomByID(message.RoomID)
	if room == nil {
		return
	}

	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client
}
func (client *Client) handleJoinRoomMessage(message MessageSocket) {
	roomID := message.RoomID
	client.joinRoom(roomID, nil)
}

func countUserChanInClient(room *RoomSocket, sender *Client) int {
	userID := sender.UserID
	count := 0
	for k, _ := range room.clients {
		if k.UserID == userID {
			count += 1
		}
	}
	return count
}

func (client *Client) joinRoom(roomID primitive.ObjectID, sender *Client) {
	room := client.wsServer.findRoomByID(roomID)
	if room != nil {
		// if a user is logged in on multiple devices (each device is chan, max=5)
		if countUserChanInClient(room, client) >= 5 {
			for k, _ := range room.clients {
				delete(room.clients, k)
				break
			}
		}
	}

	if room == nil {
		room = client.wsServer.createRoom(roomID, sender != nil)
	}

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return
	}

	if !client.isInRoom(room) {

		client.rooms[room] = true
		room.register <- client

	}

}
func (client *Client) isInRoom(room *RoomSocket) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}

	return false
}
