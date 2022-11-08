package core

type IAVCache interface {
	ReadRawFrame(frame *AVFrame) error
	ReadFrame(frame *AVFrame) error
	SeekToLastKeyFrame() error
	SeekToEarliestKeyFrame() error
}

type AVCache struct {
}

func NewAVCache() *AVCache {
	return &AVCache{}
}

func (cache *AVCache) ReadRawFrame(frame *AVFrame) error {
	return nil
}

func (cache *AVCache) ReadFrame(frame *AVFrame) error {
	return nil
}

func (cache *AVCache) SeekToLastKeyFrame() error {
	return nil
}

func (cache *AVCache) SeekToEarliestKeyFrame() error {
	return nil
}
