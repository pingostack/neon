package rtsp

import (
	"github.com/pingopenstack/neon/pkg/tcp"
	"github.com/pingopenstack/neon/src/core"
	"github.com/pingopenstack/neon/src/protocol"
	"github.com/sirupsen/logrus"
)

type Session struct {
	tcp.Context
	core.Session
	protocol.Rtsp
	logger *logrus.Entry
}

func NewSession(logger *logrus.Entry) *Session {
	conn := &Session{
		logger: logger,
	}

	return conn
}

func (s *Session) OnTcpClose() error {
	return nil
}

func (s *Session) OnTcpRread(buf []byte) (int, error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("rtsp session[%v] panic recover, %v", s, err)
		}
	}()

	if buf == nil || len(buf) == 0 {
		s.logger.Errorf("rtsp session[%v] no buf", s)
		return 0, nil
	}

	offset, err := s.Feed(buf)

	if err != nil {
		return offset, err
	}

	return offset, nil
}
