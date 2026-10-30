package types

import "encoding/json"

// message received by client peer
type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
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

// sent by a lobby to every player in it whenever its roster changes
type LobbyUpdate struct {
	LobbyID    string `json:"lobbyId"`
	Players    int    `json:"players"`
	MaxPlayers int    `json:"maxPlayers"`
}

// sent by a lobby to a player that has been removed from it
type MatchCancelled struct {
	LobbyID string `json:"lobbyId"`
}
