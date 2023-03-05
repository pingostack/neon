package forwarding

import "github.com/sirupsen/logrus"

type Session struct {
	*logrus.Entry
	context map[interface{}]interface{}
	actived bool
}

func NewSession(logger *logrus.Entry) *Session {
	return &Session{
		Entry:   logger,
		actived: false,
	}
}

func (s *Session) Active() {
	s.actived = true
}

func (s *Session) Inactive() {
	s.actived = false
}

func (s *Session) IsActive() bool {
	return s.actived
}

func (s *Session) GetLogger() *logrus.Entry {
	return s.Entry
}

func (s *Session) SetContext(k interface{}, v interface{}) {
	s.context[k] = v
}

func (s *Session) GetContext(k interface{}) interface{} {
	v, ok := s.context[k]
	if !ok {
		return nil
	}

	return v
}
