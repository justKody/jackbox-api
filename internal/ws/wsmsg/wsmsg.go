package wsmsg

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/justKody/jackbox-api/internal/ws/types"
)

// Write sends data to the client wrapped in a message of type msgType.
// A websocket allows one writer at a time, callers make sure of that
func Write(conn *websocket.Conn, msgType string, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", msgType, err)
	}
	return conn.WriteJSON(types.WSMessage{Type: msgType, Data: b})
}

// Decode reads a message's data into v, a message without data leaves v as is
// (findMatch may come without any)
func Decode(msg types.WSMessage, v any) error {
	if len(msg.Data) == 0 {
		return nil
	}
	if err := json.Unmarshal(msg.Data, v); err != nil {
		return fmt.Errorf("bad data for %s: %w", msg.Type, err)
	}
	return nil
}

// AllowOrigin reports whether a websocket upgrade from origin to host is allowed:
// no origin header, any localhost dev origin, and same host origins
// (React Native on a phone sets the origin to the server it dials)
func AllowOrigin(origin, host string) bool {
	return origin == "" || strings.HasPrefix(origin, "http://localhost") || strings.TrimPrefix(origin, "http://") == host
}
