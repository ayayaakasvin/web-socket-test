package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserInfo     *User
	Conn         *websocket.Conn
	ConnectionID string
}

func NewClient(conn *websocket.Conn, userInfo *User) *Client {
	return &Client{
		UserInfo: userInfo,
		Conn: conn,
		ConnectionID: uuid.NewString(),
	}
}