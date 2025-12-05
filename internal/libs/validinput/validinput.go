// validinput package provides custom functions to check validity of input objects like user form and etc.
package validinput

import (
	"errors"
	"slices"
	"time"
	"web-socket-test/internal/models"
	"web-socket-test/internal/models/dto"
)

// Message validation errors
const (
	EmptyPayload       = "empty payload"
	InvalidMessageType = "invalid message type"
	MessageTooBig      = "payload too large"
)

// Limits of message
const (
	// 8 KB of data
	MaxMessageSize = 8192
)

func ValidateWBSMessage(c *models.Client, msg *dto.WBSMessage) error {
	if msg.Payload == "" {
		return errors.New(EmptyPayload)
	}

	if msg.Type == "" || !slices.Contains(dto.Types, msg.Type) {
		return errors.New(InvalidMessageType)
	}

	if len(msg.Payload) > MaxMessageSize {
		return errors.New(MessageTooBig)
	}

	msg.Time = time.Now().Format(time.RFC3339)

	msg.Origin = dto.Origin{
		UserID: c.UserInfo.ID,
		Username: c.UserInfo.Username,
		System: false,
	}

	return nil
}
