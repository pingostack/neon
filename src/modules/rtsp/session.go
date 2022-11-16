package rtsp

import (
	"github.com/google/uuid"
	"github.com/pingopenstack/neon/pkg/protocol/rtsp"
	"github.com/pingopenstack/neon/pkg/tcp"
	"github.com/pingopenstack/neon/src/core"
	"github.com/sirupsen/logrus"
)

type Session struct {
	tcp.IContext
	*core.Session
	*rtsp.Serv
}

func NewSession() *Session {
	ctx := tcp.NewContext()
	logger := logrus.WithFields(logrus.Fields{
		"Module":      "RtspModule",
		"SessionType": "Rtsp",
		"SessionId":   uuid.New(),
		"RemoteAddr":  ctx.RemoteAddr(),
		"LocalAddr":   ctx.LocalAddr(),
	})

	s := &Session{
		IContext: ctx,
		Session:  core.NewSession(logger),
	}

	s.Serv = rtsp.NewServ(s.Session.Entry, func(data []byte) error {
		return s.Write(data)
	})

	s.Infof("Session created")

	return s
}

func (s *Session) OnTcpClose() error {
	s.SetConn(nil)

	return nil
}
