package forwarding

import "github.com/let-light/neon/pkg/av"

type IAVCache interface {
	ReadRawFrame(frame *av.AVFrame) error
	ReadFrame(frame *av.AVFrame) error
	SeekToLastKeyFrame() error
	SeekToEarliestKeyFrame() error
}

type AVCache struct {
}

func NewAVCache() *AVCache {
	return &AVCache{}
}

func (cache *AVCache) ReadRawFrame(frame *av.AVFrame) error {
	return nil
}

func (cache *AVCache) ReadFrame(frame *av.AVFrame) error {
	return nil
}

func (cache *AVCache) SeekToLastKeyFrame() error {
	return nil
}

func (cache *AVCache) SeekToEarliestKeyFrame() error {
	return nil
}
