package rtsp

import (
	"github.com/google/uuid"
	"github.com/let-light/network/tcp"
	"github.com/pingopenstack/neon/pkg/protocol/rtsp"
	"github.com/pingopenstack/neon/src/core"
	"github.com/sirupsen/logrus"
)

type Session struct {
	tcp.IConnection
	*core.Session
	*rtsp.Serv
}

func NewSession(c tcp.IConnection) *Session {

	logger := logrus.WithFields(logrus.Fields{
		"Module":      "RtspModule",
		"SessionType": "Rtsp",
		"SessionId":   uuid.New(),
		"RemoteAddr":  c.RemoteAddr(),
		"LocalAddr":   c.LocalAddr(),
	})

	s := &Session{
		IConnection: c,
		Session:     core.NewSession(logger),
	}

	s.Serv = rtsp.NewServ(s.Session.Entry, func(data []byte) error {
		return s.Write(data)
	})

	return s
}

func (s *Session) OnTcpClose() error {
	s.SetConn(nil)

	return nil
}
