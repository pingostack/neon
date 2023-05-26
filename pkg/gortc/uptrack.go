package gortc

import (
	"io"

	"github.com/let-light/neon/pkg/buffer"
	"github.com/let-light/neon/pkg/forwarder"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

type UpTrack struct {
	buff       *buffer.Buffer
	track      *webrtc.TrackRemote
	receiver   *webrtc.RTPReceiver
	maxBitRate uint64
	logger     *logrus.Entry
}

func NewUpTrack(buff *buffer.Buffer,
	remoteTrack *webrtc.TrackRemote,
	receiver *webrtc.RTPReceiver,
	maxBitRate uint64,
	logger *logrus.Entry) *UpTrack {
	track := &UpTrack{
		buff:       buff,
		track:      remoteTrack,
		receiver:   receiver,
		maxBitRate: maxBitRate,
		logger:     logger,
	}

	buff.OnTransportWideCC(func(sn uint16, timeNS int64, marker bool) {
		track.logger.Infof("sn %d", sn)
	})

	buff.OnFeedback(func(fb []rtcp.Packet) {

	})

	buff.Bind(receiver.GetParameters(), buffer.Options{
		MaxBitRate: maxBitRate,
	})

	go func() {
		for {
			rtpPacket, err := buff.ReadExtended()
			if err == io.EOF {
				return
			}

			logger.Debugf("read packet %+v", rtpPacket)
		}
	}()

	return track
}

func (ut *UpTrack) RequestKeyFrame() error {
	return nil
}

func (ut *UpTrack) ReadFrame() (*forwarder.FramePacket, error) {
	return nil, nil
}

func (ut *UpTrack) GetFrameKind() forwarder.FrameKind {
	return forwarder.FrameKindVideo
}

func (ut *UpTrack) GetPacketType() forwarder.PacketType {
	return forwarder.PacketTypeRtp
}

func (ut *UpTrack) GetCodecType() forwarder.CodecType {
	return forwarder.CodecTypeH264
}

func (ut *UpTrack) Close() {

}
