package types

// message received by client peer
type WSMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

// action called via player to the game server
type PlayerAction int

const (
	FindMatch PlayerAction = iota
	CancelMatch
	JoinRoom
	LeaveRoom
)

type PlayerActionMessage struct {
	Action PlayerAction
	Msg    any
}
