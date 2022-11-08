package protocol

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/pingopenstack/neon/src/core"
	"github.com/sirupsen/logrus"
)

type RtspRole int

const (
	RtspRoleServer RtspRole = iota
	RtspRoleClient
)

type State int

const (
	Empty State = iota
	Options
	Describe
	Setup
	Play
	Pause
	Teardown
)

type ConnectionMode int

const (
	ConnectionModeTcp ConnectionMode = iota
	ConnectionModeUdp
)

type IRtspImpl interface {
	RtspCmdHandler(p *Rtsp, cmdPacket *RtspCmdPacket) error
	RtpRtcpHandler(p *Rtsp, frame *core.AVFrame) error
}

type RtspCmdLines map[string]string

type RtspCmdPacket struct {
	CmdLines RtspCmdLines
	Content  []byte
}

type Rtsp struct {
	cm    ConnectionMode
	state State
	role  RtspRole // role of the Rtsp
	impl  IRtspImpl
}

func NewRtsp(role RtspRole, cm ConnectionMode, impl IRtspImpl) (*Rtsp, error) {
	if role != RtspRoleServer && role != RtspRoleClient {
		return nil, errors.New("invalid role")
	}

	return &Rtsp{
		cm:    cm,
		state: Empty,
		role:  role,
		impl:  impl,
	}, nil
}

func (p *Rtsp) State() State {
	return p.state
}

func (p *Rtsp) decodeCmd(buf []byte) (*RtspCmdPacket, int, error) {
	headerEndOffset := bytes.Index(buf, []byte("\r\n\r\n"))
	if headerEndOffset == -1 {
		return nil, -1, errors.New("incomplete packet")
	}

	logrus.Debugf("rtsp packet: headerEndOffset %d", headerEndOffset)

	contentLength := 0

	rtspPacket := &RtspCmdPacket{
		CmdLines: make(RtspCmdLines),
		Content:  nil,
	}

	lines := bytes.Split(buf, []byte("\r\n"))
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
		rtspPacket.CmdLines[key] = value

		if key == "content-length" {
			var err error
			contentLength, err = strconv.Atoi(value)
			if err != nil {
				return nil, headerEndOffset + 4, err
			}
		}
	}

	rtspPacket.Content = buf[headerEndOffset+4 : headerEndOffset+4+contentLength]

	endOffset := headerEndOffset + 4 + contentLength

	logrus.Debugf("rtsp packet: %+v", rtspPacket)

	return rtspPacket, endOffset, nil
}

func (p *Rtsp) decodeRtpRtcp(buf []byte) (*core.AVFrame, int, error) {
	return nil, 0, nil
}

func (p *Rtsp) Feed(buf []byte) (int, error) {
	logrus.Debugf("rtsp packet: %s", buf)
	if p.role == RtspRoleServer {
		return p.feedServer(buf)
	} else {
		return p.feedClient(buf)
	}
}

func (p *Rtsp) feedServer(buf []byte) (int, error) {
	var err error
	var endOffset int
	if buf[0] != '$' {
		var cmdPacket *RtspCmdPacket
		cmdPacket, endOffset, err = p.decodeCmd(buf)
		if err != nil {
			return endOffset, err
		}

		if p.impl != nil {
			err = p.impl.RtspCmdHandler(p, cmdPacket)
			if err != nil {
				return endOffset, err
			}
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

func (p *Rtsp) feedClient(buf []byte) (int, error) {
	return 0, nil
}
