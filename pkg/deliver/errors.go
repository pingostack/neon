package deliver

import "errors"

var (
	ErrFrameSourceClosed          = errors.New("frame source closed")
	ErrFrameDestinationClosed     = errors.New("frame destination closed")
	ErrFrameDestinationExists     = errors.New("frame destination exists")
	ErrCodecNotSupported          = errors.New("codec not supported")
	ErrFrameSourceAudioNotSupport = errors.New("audio not support")
	ErrFrameSourceVideoNotSupport = errors.New("video not support")
)
