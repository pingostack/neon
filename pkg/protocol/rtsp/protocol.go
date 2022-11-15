package rtsp

import (
	"errors"
	"strconv"
	"time"

	"github.com/pingopenstack/neon/src/core"
	"github.com/sirupsen/logrus"
)

type RtspRole int
type State int

type HeaderLines map[string]string

const (
	RtspRoleServer RtspRole = iota
	RtspRoleClient
)

const (
	EmptyState State = iota
	OptionsState
	DescribeState
	SetupState
	PlayState
	PauseState
	TeardownState
)

type IRtspImpl interface {
	RtspCmdHandler(p *Protocol, req *Request) error
	RtpRtcpHandler(p *Protocol, frame *core.AVFrame) error
}

type Protocol struct {
	state       State
	impl        IRtspImpl
	cseqCounter int
	settings    ProtocolSettings
}

type ProtocolSettings struct {
	Role  RtspRole
	Url   string
	Write func(date []byte) error
}

func NewProtocol(settings ProtocolSettings, impl IRtspImpl) (*Protocol, error) {
	if settings.Role != RtspRoleServer && settings.Role != RtspRoleClient {
		return nil, errors.New("invalid role")
	}

	return &Protocol{
		state:       EmptyState,
		impl:        impl,
		cseqCounter: 0,
		settings:    settings,
	}, nil
}

func (p *Protocol) State() State {
	return p.state
}

func (p *Protocol) decodeRtpRtcp(buf []byte) (*core.AVFrame, int, error) {
	return nil, 0, nil
}

func (p *Protocol) Feed(buf []byte) (int, error) {
	if p.settings.Role == RtspRoleServer {
		return p.feedServer(buf)
	} else {
		return p.feedClient(buf)
	}
}

func (p *Protocol) handleRequest(req *Request) error {
	logrus.Debugf("rtsp request: %s", req.String())
	switch req.Method() {
	case OptionMethod:
		respStr := p.NewOptionsResponse(req.CSeq(), []string{"OPTIONS", "DESCRIBE", "SETUP", "PLAY", "PAUSE", "TEARDOWN", "GET_PARAMETER"}).String()
		logrus.Debugf("rtsp response: %s", respStr)
		p.settings.Write([]byte(respStr))

	case DescribeMethod:
	case AnnounceMethod:
	case SetupMethod:
	case PlayMethod:
	case PauseMethod:
	case TeardownMethod:
	case GetParameterMethod:
	case SetParameterMethod:
	case RecordMethod:
	default:
	}

	return nil
}

func (p *Protocol) feedServer(buf []byte) (int, error) {
	var err error
	var endOffset int
	if buf[0] != '$' {
		var req *Request
		req, endOffset, err = UnmarshalRequest(buf)
		if err != nil {
			return endOffset, err
		}

		err := p.handleRequest(req)
		if err != nil {
			return endOffset, err
		}

		// if p.impl != nil {
		// 	err = p.impl.RtspCmdHandler(p, req)
		// 	if err != nil {
		// 		return endOffset, err
		// 	}
		// }

	} else {
		var frame *core.AVFrame
		frame, endOffset, err = p.decodeRtpRtcp(buf)
		if err != nil {
			return endOffset, err
		}

		if p.impl != nil {
			err = p.impl.RtpRtcpHandler(p, frame)
			if err != nil {
				return endOffset, err
			}
		}
	}

	return endOffset, nil
}

func (p *Protocol) feedClient(buf []byte) (int, error) {
	var err error
	var endOffset int
	if buf[0] != '$' {
		_, endOffset, err = UnmarshalResponse(buf)
		if err != nil {
			return endOffset, err
		}

		// err = p.responseProcess(resp)
		// if err != nil {
		// 	return endOffset, err
		// }
	} else {
		var frame *core.AVFrame
		frame, endOffset, err = p.decodeRtpRtcp(buf)
		if err != nil {
			return endOffset, err
		}

		if p.impl != nil {
			err = p.impl.RtpRtcpHandler(p, frame)
			if err != nil {
				return endOffset, err
			}
		}
	}

	return endOffset, nil
}

func (p *Protocol) nextCSeq() string {
	p.cseqCounter++
	return strconv.Itoa(p.cseqCounter)
}

func (p *Protocol) NewRequest(method string) *Request {
	req := &Request{
		method:  method,
		url:     "*",
		version: "RTSP/1.0",
		lines: HeaderLines{
			"CSeq": p.nextCSeq(),
		},
	}

	return req
}

func (p *Protocol) NewOptionsRequest() *OptionsRequest {
	req := p.NewRequest("OPTIONS").Option()
	return req
}

func (p *Protocol) NewDescribeRequest(sdp string) *DescribeRequest {
	req := p.NewRequest("DESCRIBE").Describe()
	req.SetSdp(sdp)
	return req
}

func (p *Protocol) NewAnnounceRequest() *AnnounceRequest {
	req := p.NewRequest("ANNOUNCE").Announce()
	return req
}

func (p *Protocol) NewSetupRequest(trackId int, transport *Transport) *SetupRequest {
	req := p.NewRequest("SETUP").Setup()
	req.SetTransport(transport)

	return req
}

func (p *Protocol) NewPlayRequest() *PlayRequest {
	req := p.NewRequest("PLAY").Play()
	return req
}

func (p *Protocol) NewPauseRequest() *PauseRequest {
	req := p.NewRequest("PAUSE").Pause()
	return req
}

func (p *Protocol) NewTeardownRequest() *TeardownRequest {
	req := p.NewRequest("TEARDOWN").Teardown()
	return req
}

func (p *Protocol) NewGetParameterRequest() *GetParameterRequest {
	req := p.NewRequest("GET_PARAMETER").GetParameter()
	return req
}

func (p *Protocol) NewSetParameterRequest() *SetParameterRequest {
	req := p.NewRequest("SET_PARAMETER").SetParameter()
	return req
}

func (p *Protocol) NewRecordRequest() *RecordRequest {
	req := p.NewRequest("RECORD").Record()
	return req
}

func (p *Protocol) NewResponse(cseq int, status Status) *Response {
	resp := &Response{
		version:   "RTSP/1.0",
		status:    status,
		statusStr: status.String(),
		lines: HeaderLines{
			"CSeq": strconv.Itoa(cseq),
			"Date": time.Now().Format(time.RFC1123),
		},
	}

	return resp
}

func (p *Protocol) NewOptionsResponse(cseq int, options []string) *OptionsResponse {
	resp := p.NewResponse(cseq, StatusOK).Option()
	resp.SetOptions(options)
	return resp
}
