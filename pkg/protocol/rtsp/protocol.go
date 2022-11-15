package rtsp

import (
	"errors"

	"github.com/pingopenstack/neon/src/core"
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
	state State
	role  RtspRole // role of the Protocol
	impl  IRtspImpl
}

func NewProtocol(role RtspRole, impl IRtspImpl) (*Protocol, error) {
	if role != RtspRoleServer && role != RtspRoleClient {
		return nil, errors.New("invalid role")
	}

	return &Protocol{
		state: EmptyState,
		role:  role,
		impl:  impl,
	}, nil
}

func (p *Protocol) State() State {
	return p.state
}

func (p *Protocol) decodeRtpRtcp(buf []byte) (*core.AVFrame, int, error) {
	return nil, 0, nil
}

func (p *Protocol) Feed(buf []byte) (int, error) {
	if p.role == RtspRoleServer {
		return p.feedServer(buf)
	} else {
		return p.feedClient(buf)
	}
}

func (p *Protocol) requestProcess(req *Request) error {
	return nil
}

func (p *Protocol) responseProcess(resp *Response) error {
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

		err = p.requestProcess(req)
		if err != nil {
			return endOffset, err
		}

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
		var resp *Response
		resp, endOffset, err = UnmarshalResponse(buf)
		if err != nil {
			return endOffset, err
		}

		err = p.responseProcess(resp)
		if err != nil {
			return endOffset, err
		}
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
