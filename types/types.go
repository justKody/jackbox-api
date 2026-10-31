package types

import (
	"encoding/json"

	"github.com/justKody/jackbox-api/game"
	"github.com/justKody/jackbox-api/game/rank"
	gametype "github.com/justKody/jackbox-api/game/types"
)

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
	LockGuess
	Vote
	Next
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

// sent by the client with findMatch, who the player is and what they look like
type Profile struct {
	// stable id the app keeps on device, ratings are keyed by it
	ClientID string `json:"clientId"`
	Name     string `json:"name"`
	Animal   string `json:"animal"`
	Skin     string `json:"skin"`
	// the RP the app has saved, only used the first time the server sees this client
	RP int `json:"rp"`
}

// data of a lockGuess message
type LockGuessMsg struct {
	Guess int `json:"guess"`
}

// data of a dial message, the guesser turning the dial before locking in
type DialMsg struct {
	Value int `json:"value"`
}

// sent to everyone but the guesser as the guesser turns the dial
type DialUpdate struct {
	PlayerID string `json:"playerId"`
	Value    int    `json:"value"`
}

// data of a vote message
type VoteMsg struct {
	Call gametype.Call `json:"call"`
}

// sent by a lobby to every player in it whenever its roster changes
type LobbyUpdate struct {
	LobbyID    string            `json:"lobbyId"`
	You        string            `json:"you"`
	Players    []gametype.Player `json:"players"`
	MaxPlayers int               `json:"maxPlayers"`
	// seconds until the game starts, 0 while waiting for more players
	StartsIn int `json:"startsIn"`
}

// sent by a lobby to a player that has been removed from it
type MatchCancelled struct {
	LobbyID string `json:"lobbyId"`
}

// sent by a lobby to every player whenever the game changes
type GameUpdate struct {
	// the receiving player's id in the game
	You string `json:"you"`
	game.View
}

// sent to each player when a ranked game ends
type RankResult = rank.Result

// sent to a player when one of their actions was rejected
type ActionError struct {
	Message string `json:"message"`
}
