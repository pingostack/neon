package rtc

type Kind int

const (
	KindAudio Kind = iota
	KindVideo
)

type Track struct {
	kind          Kind
	rtpParameters *RtpParameters
}

func (t *Track) Kind() Kind {
	return t.kind
}

func (t *Track) RtpParameters() *RtpParameters {
	return t.rtpParameters
}

func (t *Track) SetRtpParameters(rtpParameters *RtpParameters) {
	t.rtpParameters = rtpParameters
}
