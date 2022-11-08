package core

type PayloadType int

const (
	AudioPayloadTypeAAC PayloadType = iota // 0
	AudioPayloadTypeOpus
	AudioPayloadTypeNum
	VideoPayloadTypeH264
	VideoPayloadTypeH265
	VideoPayloadTypeVP8
	VideoPayloadTypeVP9
	VideoPayloadTypeAV1
	VideoPayloadTypeNum
)

type PacketType int

const (
	PacketRaw PacketType = iota // 0
	PacketRtmp
	PacketMpegts
	PacketRtp
	PacketRtcp
)

type AVFrame struct {
	PT      PayloadType
	PktType PacketType
	PTS     uint64
	DTS     uint64
	Payload []byte
}

func (frame *AVFrame) IsVideo() bool {
	return frame.PT > AudioPayloadTypeNum && frame.PT < VideoPayloadTypeNum
}

func (frame *AVFrame) IsAudio() bool {
	return frame.PT < AudioPayloadTypeNum
}
