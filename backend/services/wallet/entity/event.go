package entity

// Event defines event.
type Event struct {
	Topic   string
	Key     []byte
	Payload []byte
}
