package types

type WSMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}
