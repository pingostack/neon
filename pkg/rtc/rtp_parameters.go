package rtc

func (t RtcType) String() string {
	switch t {
	case RtcTypePipe:
		return "pipe"
	case RtcTypeSimple:
		return "simple"
	case RtcTypeSvc:
		return "svc"
	case RtcTypeSimulcast:
		return "simulcast"
	default:
		return "none"
	}
}

func String2RtcType(t string) RtcType {
	switch t {
	case "pipe":
		return RtcTypePipe
	case "simple":
		return RtcTypeSimple
	case "svc":
		return RtcTypeSvc
	case "simulcast":
		return RtcTypeSimulcast
	default:
		return RtcTypeNone
	}
}

func (rtpParameters *RtpParameters) GetType() RtcType {
	if len(rtpParameters.Encodings) == 1 {
		encoding := rtpParameters.Encodings[0]
		if encoding.SpatialLayers > 1 || encoding.TemporalLayers > 1 {
			return RtcTypeSvc
		} else {
			return RtcTypeSimple
		}
	} else if len(rtpParameters.Encodings) > 1 {
		return RtcTypeSimulcast
	} else {
		return RtcTypeNone
	}
}

func (rtpParameters *RtpParameters) GetCodec(payloadType uint8) *RtpCodecParameters {
	for _, codec := range rtpParameters.Codecs {
		if codec.PayloadType == payloadType {
			return &codec
		}
	}
	return nil
}

func (rtpParameters *RtpParameters) GetRtxCodec(payloadType uint8) *RtpCodecParameters {
	for _, codec := range rtpParameters.Codecs {
		apt, err := codec.Parameters.Uint8("apt")
		if err == nil && apt == payloadType {
			return &codec
		}
	}
	return nil
}
