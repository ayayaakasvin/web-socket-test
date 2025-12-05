package postgresql

import (
	"context"
	"web-socket-test/internal/models/dto"
)

// SaveMessage implements [core.ChatHistoryStorage].
func (p *PostgreSQL) SaveMessage(ctx context.Context, wm *dto.WBSMessage) error {
	_, err := p.conn.ExecContext(
		ctx,
		`INSERT INTO chat_history (payload, user_id, type, created_at)
		 VALUES ($1, $2, $3, $4)`,
		wm.Payload,
		wm.Origin.UserID,
		wm.Type,
		wm.Time,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetRecentMessage implements [core.ChatHistoryStorage].
func (p *PostgreSQL) GetRecentMessage(ctx context.Context, limit int) ([]*dto.WBSMessage, error) {
	rows, err := p.conn.QueryContext(
		ctx,
		`SELECT h.message_id, h.payload, h.type, h.created_at, u.username
		FROM chat_history h
		LEFT JOIN users u ON h.user_id = u.user_id
		ORDER BY h.message_id DESC 
		LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}

	var msgs []*dto.WBSMessage
	for rows.Next() {
		msg := new(dto.WBSMessage)

        if err := rows.Scan(&msg.ID, &msg.Payload, &msg.Type, &msg.Time, &msg.Origin.UserID, &msg.Origin.Username); err != nil {
            return msgs, err
        }

		switch msg.Type {
		case dto.ConnectType, dto.DisconnectType:
			msg.Origin.System = true
		default:
			msg.Origin.System = false
		}

		msgs = append(msgs, msg)
	}

	if rows.Err() != nil {
		return msgs, rows.Err()
	}

	return msgs, nil
}
