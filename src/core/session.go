package core

type Session struct {
	context map[interface{}]interface{}
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
