package dto

import "time"

type WBSMessage struct {
	ID      uint   `json:"message_id"` // message id
	Type    string `json:"type"`       // "message", "connect", "disconnect", etc.
	Time    string `json:"time"`       // timestamp
	Origin  Origin `json:"origin"`     // user/system info
	Payload string `json:"payload"`    // actual message text (or empty for system events)
}

type Origin struct {
	UserID   uint   `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	System   bool   `json:"system"`
}

const (
	ConnectType    = "connect"
	MessageType    = "message"
	DisconnectType = "disconnect"
)

var Types []string = []string{ConnectType, MessageType, DisconnectType}

func SystemMessage(userID uint, username string, action string) *WBSMessage {
	payload := ""
	switch action {
	case ConnectType:
		payload = ConnectType
	case DisconnectType:
		payload = DisconnectType
	}

	og := Origin{
		System: true,
		UserID: userID,
		Username: username,
	}

	return &WBSMessage{
		Type: action,
		Time: time.Now().Format(time.RFC3339),
		Origin: og,
		Payload: payload,
	}
}