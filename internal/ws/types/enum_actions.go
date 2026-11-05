package types

// action called via player to the game server
type PlayerAction int

const (
	FindMatch PlayerAction = iota
	CancelMatch
	JoinRoom
	LeaveRoom
	LockGuess
	Vote
	Dial
)

type PlayerActionMessage struct {
	Action PlayerAction
	Msg    any
}

// action the game server sends to itself, e.g. from a timer
type ServerAction int

const (
	PlacePlayer ServerAction = iota
)

type ServerActionMessage struct {
	Action ServerAction
	Msg    any
}
