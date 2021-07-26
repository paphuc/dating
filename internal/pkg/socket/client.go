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
	wsServer *WsServer
	conn     *websocket.Conn
	send     chan *MessageSocket
	save     *chan SaveMessage
	rooms    map[*RoomSocket]bool
}

func NewClient(conn *websocket.Conn, wsServer *WsServer, idRoom primitive.ObjectID, sm *chan SaveMessage) *Client {
	return &Client{
		ID:       primitive.NewObjectID(),
		RoomId:   idRoom,
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
			logger.Errorf("Failed when read message from room %d, client %d", c.RoomId, c.ID, mgs)
			return
		}
		c.handleNewMessage(mgs)
	}
}

func (c *Client) Write(logger glog.Logger) {
	defer c.conn.Close()
	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			logger.Errorf("Failed when write message from room %d, client %d", c.RoomId, c.ID, msg)
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
			sm := &SaveMessage{
				message: &types.Message{
					ID:          jsonMessage.ID,
					RoomID:      roomID,
					SenderID:    jsonMessage.SenderID,
					Content:     jsonMessage.Content,
					Attachments: jsonMessage.Attachments,
					CreateAt:    jsonMessage.CreateAt,
				},
			}
			*client.save <- *sm
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

func (client *Client) joinRoom(roomID primitive.ObjectID, sender *Client) {

	room := client.wsServer.findRoomByID(roomID)
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
