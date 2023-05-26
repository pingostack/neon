package forwarder

type IUpTrack interface {
	IFrameSource
	StreamID() string
	TrackID() string
	Layer() int
}
