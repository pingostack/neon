package protocol

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

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

func (p *Protocol) requestUnmarshal(buf []byte) (*Request, int, error) {
	headerEndOffset := bytes.Index(buf, []byte("\r\n\r\n"))
	if headerEndOffset == -1 {
		return nil, -1, errors.New("incomplete packet")
	}

	endOffset := headerEndOffset + 4

	contentLength := 0

	req := &Request{
		Lines: make(HeaderLines),
	}

	lines := bytes.Split(buf[:headerEndOffset], []byte("\r\n"))
	if len(lines) < 2 {
		return nil, endOffset, errors.New("invalid packet")
	}

	// parse first line
	methodLine := lines[0]
	methodLineParts := bytes.Split(methodLine, []byte(" "))
	if len(methodLineParts) != 3 {
		return nil, endOffset, errors.New("invalid packet")
	}

	req.Method = strings.ToLower(string(methodLineParts[0]))
	req.Url = strings.ToLower(string(methodLineParts[1]))
	req.Version = strings.ToLower(string(methodLineParts[2]))

	// parse other lines
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		idx := bytes.Index(line, []byte(":"))
		if idx == -1 {
			continue
		}

		key := strings.ToLower(string(line[:idx]))
		value := string(line[idx+1:])
		req.Lines[key] = value

		if key == "content-length" {
			var err error
			contentLength, err = strconv.Atoi(value)
			if err != nil {
				return nil, headerEndOffset + 4, err
			}
		}
	}

	req.Content = buf[headerEndOffset+4 : headerEndOffset+4+contentLength]

	endOffset += contentLength

	return req, endOffset, nil
}

func (p *Protocol) feedServer(buf []byte) (int, error) {
	var err error
	var endOffset int
	if buf[0] != '$' {
		var req *Request
		req, endOffset, err = p.requestUnmarshal(buf)
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
	return 0, nil
}
