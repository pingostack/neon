package rtsp

type RtspRole int
type State int
type HeaderLines map[string]string
type WriteHandler func(date []byte) error
