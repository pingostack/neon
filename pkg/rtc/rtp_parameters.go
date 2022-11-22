package rtc

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
