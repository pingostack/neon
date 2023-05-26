package forwarder

type FrameKind int

const (
	FrameKindVideo FrameKind = 1
	FrameKindAudio           = 2
	FrameKindData            = 4
)

type PacketType int16

const (
	PacketTypeRtp PacketType = iota + 1
	PacketTypeRtcp
	PacketTypeRtmp
	PacketTypeFLV
)

type CodecType int16

const (
	CodecTypeVP8 CodecType = iota + 1
	CodecTypeVP9
	CodecTypeH264
	CodecTypeVP8SVC
	CodecTypeVP9SVC
	CodecTypeH264SVC
	CodecTypeH265
	CodecTypeAV1
	CodecTypeAV1SVC
	CodecTypeOpus
	CodecTypeG722
	CodecTypePCMU
	CodecTypePCMA
	CodecTypeAAC
)

type FramePacket struct {
	Kind     FrameKind
	PT       PacketType
	Codec    CodecType
	Layer    int8
	KeyFrame bool
	Payload  []byte
}

type IFrameSource interface {
	RequestKeyFrame() error
	ReadFrame() (*FramePacket, error)
	GetFrameKind() FrameKind
	GetPacketType() PacketType
	GetCodecType() CodecType
	Close()
}

type IFrameDestination interface {
	WriteFrame(frame *FramePacket) error
	Close()
}

type FrameSource struct {
	kind  FrameKind
	pt    PacketType
	codec CodecType
}

func NewDefaultSource(kind FrameKind, pt PacketType, codec CodecType) IFrameSource {
	return &FrameSource{
		kind:  kind,
		pt:    pt,
		codec: codec,
	}
}

func (fs *FrameSource) RequestKeyFrame() error {
	return nil
}

func (fs *FrameSource) ReadFrame() (*FramePacket, error) {
	return nil, nil
}

func (fs *FrameSource) GetCodecType() CodecType {
	return fs.codec
}

func (fs *FrameSource) GetFrameKind() FrameKind {
	return fs.kind
}

func (fs *FrameSource) GetPacketType() PacketType {
	return fs.pt
}

func (fs *FrameSource) Close() {
}
