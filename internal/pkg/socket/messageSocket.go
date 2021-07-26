package socket

import "dating/internal/app/api/types"

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"

type MessageSocket struct {
	Action string `json:"action"`
	types.Message
}
