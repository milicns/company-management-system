package utils

type Event struct {
	Topic     string
	EventType string
	Value     string
}

func NewEvent(value, eventType string) Event {
	return Event{
		Topic:     eventType,
		EventType: eventType,
		Value:     value,
	}
}
