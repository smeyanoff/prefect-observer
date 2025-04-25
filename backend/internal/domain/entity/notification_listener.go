package entity

type Listener interface {
	WriteMessage(messageType int, data []byte) error
	Close() error
}
