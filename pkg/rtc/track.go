package rtc

type Kind int

const (
	KindAudio Kind = iota
	KindVideo
)

type Track struct {
	kind          Kind
	rtpParameters RtpParameters
	//struct RtpMapping rtpMapping;
}
