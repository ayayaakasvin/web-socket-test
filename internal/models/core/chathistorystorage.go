package core

type ChatHistoryStorage interface {
	MessageReader
	MessageWriter

	Close() error
}