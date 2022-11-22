package rtc

type RtcType int

const (
	RtcTypeNone RtcType = iota - 1
	RtcTypePipe
	RtcTypeSimple
	RtcTypeSvc
	RtcTypeSimulcast
)

type RtpParameters struct {
	// MID RTP MID (RFC 5888).
	MID string `json:"mid,omitempty"`

	// Codecs RTP codecs supported by this endpoint.
	Codecs []RtpCodecParameters `json:"codecs,omitempty"`

	// HeaderExtensions RTP header extensions supported by this endpoint.
	HeaderExtensions []RtpHeaderExtensionParameters `json:"headerExtensions,omitempty"`

	// Encodings RTP encodings supported by this endpoint.
	Encodings []RtpEncodingParameters `json:"encodings,omitempty"`

	// RTCP RTCP options.
	RTCP RtcpParameters `json:"rtcp,omitempty"`

	// HasRTCP Whether RTCP is enabled.
	HasRTCP bool `json:"hasRtcp,omitempty"`
}

type RtpCodecParameters struct {
	// MIME media codec MIME type/subtype.
	MIME string `json:"mimeType,omitempty"`

	// PayloadType RTP payload type.
	PayloadType uint8 `json:"payloadType,omitempty"`

	// ClockRate RTP clock rate.
	ClockRate uint32 `json:"clockRate,omitempty"`

	// Channels number of channels for audio codecs.
	Channels uint8 `json:"channels,omitempty"`

	// Parameters codec-specific parameters.
	Parameters map[string]string `json:"parameters,omitempty"`

	// RTCPFeedback codec-specific RTCP feedback mechanisms.
	RTCPFeedback []RtcpFeedback `json:"rtcpFeedback,omitempty"`
}

type MimeType uint16

const (
	MimeTypeUnknown MimeType = iota
	MimeTypeAudio
	MimeTypeVideo
)

type MineSubtype uint16

const (
	Unknown MineSubtype = 0

	// Audio
	OPUS      = 100
	MULTIOPUS = 101
	PCMA      = 102
	PCMU      = 103
	ISAC      = 104
	G722      = 105
	ILBC      = 106
	SILK      = 107

	// Video
	VP8 = 200
	VP9
	H264
	H264_SVC
	X_H264UC
	H265

	// Complementary codecs:
	CN = 300
	TELEPHONE_EVENT

	// Feature codecs:
	RTX = 400
	ULPFEC
	X_ULPFECUC
	FLEXFEC
	RED
)

type RtpCodecMimeType struct {
	MimeType MimeType
	Subtype  MineSubtype
}

type RtcpFeedback struct {
	// Type RTCP feedback type.
	Type string `json:"type,omitempty"`

	// Parameter RTCP feedback parameter.
	Parameter string `json:"parameter,omitempty"`
}

type RtpHeaderExtensionUriType int

const (
	RtpHeaderExtensionUriUnknown RtpHeaderExtensionUriType = iota
	RtpHeaderExtensionUriMid
	RtpHeaderExtensionUriRtpStreamId
	RtpHeaderExtensionUriRepariedRtpStreamId
	RtpHeaderExtensionUriAbsSendTime
	RtpHeaderExtensionUriTWCC01
	RtpHeaderExtensionUriFrameMarking07
	RtpHeaderExtensionUriFrameMarking
	RtpHeaderExtensionUriSSRCAudioLevel
	RtpHeaderExtensionUriVideoOrientation
	RtpHeaderExtensionUriTOffset
	RtpHeaderExtensionUriAbsCaptureTime
)

type RtpHeaderExtensionParameters struct {
	// URI URI of the RTP header extension.
	URI string `json:"uri,omitempty"`

	// Type RTP header extension type.
	Type RtpHeaderExtensionUriType `json:"type,omitempty"`

	// ID RTP header extension ID.
	ID uint8 `json:"id,omitempty"`

	// Encrypt whether the RTP header extension must be encrypted or not.
	Encrypt bool `json:"encrypt,omitempty"`

	// Parameters extension-specific parameters.
	Parameters map[string]string `json:"parameters,omitempty"`
}

type RtpEncodingParameters struct {
	// SSRC RTP SSRC.
	SSRC uint32 `json:"ssrc,omitempty"`
	Rid  string `json:"rid,omitempty"`

	// CodecPayloadType RTP payload type of the codec that this encoding is compatible with.
	CodecPayloadType uint8 `json:"codecPayloadType,omitempty"`

	// RTX SSRC of the retransmission RTP stream.
	RTX *RtpRtxParameters `json:"rtx,omitempty"`

	// MaxBitrate maximum bitrate in bps.
	MaxBitrate uint32 `json:"maxBitrate,omitempty"`

	// MaxFramerate maximum framerate in fps.
	MaxFramerate float64 `json:"maxFramerate,omitempty"`

	// DTX whether DTX is enabled or not.
	DTX bool `json:"dtx,omitempty"`

	// ScalabilityMode scalability mode.
	ScalabilityMode string `json:"scalabilityMode,omitempty"`

	// SpatialLayers number of spatial layers.
	SpatialLayers uint8 `json:"spatialLayers,omitempty"`

	// TemporalLayers number of temporal layers.
	TemporalLayers uint8 `json:"temporalLayers,omitempty"`

	// KSVS whether SVC scalability mode is KSVS or not.
	KSVS bool `json:"ksvs,omitempty"`

	// FEC SSRC of the RTP stream used for FEC.
	FEC *RtpFecParameters `json:"fec,omitempty"`

	// Active whether this encoding is active or not.
	Active bool `json:"active,omitempty"`

	// // ScaleResolutionDownBy scale resolution down by.
	// ScaleResolutionDownBy float64 `json:"scaleResolutionDownBy,omitempty"`

	// // Rids list of RIDs.
	// Rids []RtpRidParameters `json:"rids,omitempty"`

	// // DependencyEncodingIds list of dependency encoding IDs.
	// DependencyEncodingIds []string `json:"dependencyEncodingIds,omitempty"`

	// // Priority encoding priority.
	// Priority float64 `json:"priority,omitempty"`
}

type RtpRtxParameters struct {
	// SSRC RTP SSRC.
	SSRC uint32 `json:"ssrc,omitempty"`

	// PayloadType RTP payload type of the codec that this RTX encoding is compatible with.
	PayloadType uint8 `json:"payloadType,omitempty"`
}

type RtpFecParameters struct {
	// SSRC RTP SSRC.
	SSRC uint32 `json:"ssrc,omitempty"`

	// Mechanism FEC mechanism.
	Mechanism string `json:"mechanism,omitempty"`

	// Parameters FEC-specific parameters.
	Parameters map[string]string `json:"parameters,omitempty"`

	// RepairSSRC RTP SSRC of the RTP stream used for repair.
	RepairSSRC uint32 `json:"repairSSRC,omitempty"`
}

type RtpRidParameters struct {
	// ID RID value.
	ID string `json:"id,omitempty"`

	// Direction RID direction.
	Direction string `json:"direction,omitempty"`

	// Dependencies list of RID dependencies.
	Dependencies []string `json:"dependencies,omitempty"`

	// EncodingId encoding ID.
	EncodingId string `json:"encodingId,omitempty"`
}

type RtcpParameters struct {
	// CNAME CNAME.
	CNAME string `json:"cname,omitempty"`

	// ReducedSize whether reduced size RTCP is enabled or not.
	ReducedSize bool `json:"reducedSize,omitempty"`

	// Compound whether compound RTCP is enabled or not.
	Compound bool `json:"compound,omitempty"`

	// Mux whether RTCP-mux is enabled or not.
	Mux bool `json:"mux,omitempty"`
}

// type RtpCapabilities struct {
// 	// Codecs RTP codecs supported by this endpoint.
// 	Codecs []RtpCodecCapability `json:"codecs,omitempty"`

// 	// HeaderExtensions RTP header extensions supported by this endpoint.
// 	HeaderExtensions []RtpHeaderExtensionCapability `json:"headerExtensions,omitempty"`
// }

// type RtpCodecCapability struct {
// 	// MIME media codec MIME type/subtype.
// 	MIME string `json:"mimeType,omitempty"`

// 	// ClockRate RTP clock rate.
// 	ClockRate uint32 `json:"clockRate,omitempty"`

// 	// Channels number of channels for audio codecs.
// 	Channels uint8 `json:"channels,omitempty"`

// 	// Parameters codec-specific parameters.
// 	Parameters map[string]string `json:"parameters,omitempty"`
// }

// type RtpHeaderExtensionCapability struct {
// 	// URI URI of the RTP header extension.
// 	URI string `json:"uri,omitempty"`

// 	// PreferredId preferred RTP header extension ID.
// 	PreferredId uint8 `json:"preferredId,omitempty"`

// 	// PreferredEncrypt whether the RTP header extension must be encrypted or not.
// 	PreferredEncrypt bool `json:"preferredEncrypt,omitempty"`

// 	// Direction RTP header extension direction.
// 	Direction string `json:"direction,omitempty"`
// }

// type RtpHeaderExtensionDirection string

// const (
// 	RtpHeaderExtensionDirectionSendrecv RtpHeaderExtensionDirection = "sendrecv"
// 	RtpHeaderExtensionDirectionSendonly RtpHeaderExtensionDirection = "sendonly"
// 	RtpHeaderExtensionDirectionRecvonly RtpHeaderExtensionDirection = "recvonly"
// 	RtpHeaderExtensionDirectionInactive RtpHeaderExtensionDirection = "inactive"
// )
