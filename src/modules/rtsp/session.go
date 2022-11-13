package rtsp

import (
	"github.com/google/uuid"
	"github.com/pingopenstack/neon/pkg/tcp"
	"github.com/pingopenstack/neon/src/core"
	"github.com/pingopenstack/neon/src/modules/rtsp/protocol"
	"github.com/sirupsen/logrus"
)

type Session struct {
	tcp.IContext
	core.Session
	*protocol.Protocol
	*logrus.Entry
}

func NewSession(c interface{}, role protocol.RtspRole) *Session {
	ctx := tcp.NewContext(c)
	s := &Session{
		IContext: ctx,
		Entry: logrus.WithFields(logrus.Fields{
			"Module":      "RtspModule",
			"SessionType": "Rtsp",
			"SessionId":   uuid.New(),
			"RemoteAddr":  ctx.RemoteAddr(),
			"LocalAddr":   ctx.LocalAddr(),
		}),
	}

	var err error
	s.Protocol, err = protocol.NewProtocol(role, s)
	if err != nil {
		s.Errorf("error creating rtsp: %v", err)
		return nil
	}

	s.Infof("Session created")

	return s
}

func (s *Session) OnTcpClose() error {
	s.SetConn(nil)

	return nil
}

func (s *Session) RtspCmdHandler(p *protocol.Protocol, req *protocol.Request) error {
	return nil
}

func (s *Session) RtpRtcpHandler(p *protocol.Protocol, frame *core.AVFrame) error {
	return nil
}
