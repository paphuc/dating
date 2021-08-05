package socket

import "dating/internal/app/api/types"

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"

var (
	DoneMessage  = "done"
	ErrorMessage = "error"
)

type MessageSocket struct {
	Action string `json:"action"`
	Status string `json:"status"`
	types.Message
}
