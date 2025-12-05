package postgresql

import (
	"context"
	"time"

	"github.com/ayayaakasvin/web-socket-test/internal/models"
	"github.com/ayayaakasvin/web-socket-test/internal/models/dto"
)

type PostgreSQL_Mock struct{}

var MockUser *models.User = &models.User{
	ID:        67,
	Username:  "John Doe",
	Role:      models.AdminRole,
	CreatedAt: time.Now(),
}

var MockMessageList []*dto.WBSMessage = []*dto.WBSMessage{
	{
		ID:   1,
		Type: "connect",
		Time: time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		Origin: dto.Origin{
			UserID:   MockUser.ID,
			Username: MockUser.Username,
			System:   true,
		},
		Payload: "connected",
	},
	{
		ID:   2,
		Type: "message",
		Time: time.Now().Add(-2 * time.Minute).Format(time.RFC3339),
		Origin: dto.Origin{
			UserID:   MockUser.ID,
			Username: MockUser.Username,
			System:   false,
		},
		Payload: "Hello, this is a mock message",
	},
	{
		ID:   3,
		Type: "disconnect",
		Time: time.Now().Add(-1 * time.Minute).Format(time.RFC3339),
		Origin: dto.Origin{
			UserID:   MockUser.ID,
			Username: MockUser.Username,
			System:   true,
		},
		Payload: "user disconnected",
	},
}

// GetPublicUserInfo implements [core.UserRepository].
func (p *PostgreSQL_Mock) GetPublicUserInfo(ctx context.Context, userID uint) (*models.User, error) {
	castedMockUser := MockUser
	castedMockUser.Role = ""

	return castedMockUser, nil
}

func (p *PostgreSQL_Mock) GetPrivateUserInfo(ctx context.Context, userID uint) (*models.User, error) {
	return MockUser, nil
}

// GetRecentMessage implements [core.ChatHistoryStorage].
func (p *PostgreSQL_Mock) GetRecentMessage(ctx context.Context, limit int) ([]*dto.WBSMessage, error) {
	if limit <= 0 {
		return []*dto.WBSMessage{}, nil
	}

	total := len(MockMessageList)
	if total == 0 {
		return []*dto.WBSMessage{}, nil
	}

	if limit >= total {
		// return a copy of the slice
		res := make([]*dto.WBSMessage, total)
		copy(res, MockMessageList)
		return res, nil
	}

	// return the most recent `limit` messages
	start := total - limit
	res := make([]*dto.WBSMessage, limit)
	copy(res, MockMessageList[start:])
	return res, nil
}

// SaveMessage implements [core.ChatHistoryStorage].
func (p *PostgreSQL_Mock) SaveMessage(ctx context.Context, wm *dto.WBSMessage) error {
	// assign an ID if missing
	var maxID uint = 0
	for _, m := range MockMessageList {
		if m != nil && m.ID > maxID {
			maxID = m.ID
		}
	}
	if wm.ID == 0 {
		wm.ID = maxID + 1
	}
	if wm.Time == "" {
		wm.Time = time.Now().Format(time.RFC3339)
	}
	MockMessageList = append(MockMessageList, wm)
	return nil
}

// Close implements [core.ChatHistoryStorage].
func (p *PostgreSQL_Mock) Close() error {
	return nil
}

func NewPostgreSQL_Mock() *PostgreSQL_Mock {
	return &PostgreSQL_Mock{}
}
