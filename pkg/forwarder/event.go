package forwarder

type EventType int

const (
	EventInvalid EventType = 0
	EventJoin
	EventAudioFrame
	EventVideoFrame
)

type Event struct {
	Type EventType
}
