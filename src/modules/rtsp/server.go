package rtsp

import (
	"github.com/panjf2000/gnet"
	"github.com/pingopenstack/neon/pkg/tcp"
	"github.com/sirupsen/logrus"
)

type RtspServerSettings struct {
	Addr        string       `mapstructure:"addr"`
	TcpSettings gnet.Options `mapstructure:"tcp"`
}

type RtspServer struct {
	tcpServer *tcp.Server
	settings  RtspServerSettings
	*logrus.Entry
}

func NewRtspServer(settings RtspServerSettings) (*RtspServer, error) {
	var err error
	server := &RtspServer{
		Entry: logrus.WithFields(logrus.Fields{
			"Server": "RtspServer",
		}),
	}

	server.settings = settings

	settings.TcpSettings.Logger = logrus.StandardLogger()
	settings.TcpSettings.Codec = server

	server.tcpServer, err = tcp.NewServer(settings.Addr, server, gnet.WithOptions(settings.TcpSettings))
	if err != nil {
		server.Errorf("error creating tcp server: %v", err)
		return nil, err
	}

	return server, nil
}

func (server *RtspServer) NewOrGet(c interface{}) tcp.IContext {
	session := NewSession(c, RtspRoleServer)

	return session
}

func (server *RtspServer) OnTcpClose(ctx tcp.IContext) error {
	session := ctx.(*Session)

	return session.OnTcpClose()
}

func (server *RtspServer) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	return nil, nil
}

// Decode decodes frames from TCP stream via specific implementation.
//
// Note that when there is an incomplete packet, you should return (nil, ErrIncompletePacket)
// to make gnet continue reading data from socket, otherwise, any error other than ErrIncompletePacket
// will cause the connection to close.
func (server *RtspServer) Decode(c gnet.Conn) ([]byte, error) {
	cctx := c.Context()
	if cctx == nil {
		server.Errorf("connection context is nil")
		return nil, nil
	}

	session := cctx.(*Session)

	ctx := rtspContext(session)

	if ctx == nil {
		session.Errorf("error getting rtsp context")
		return nil, nil
	}

	offset, err := session.Feed(c.Read())
	if err != nil {
		session.Errorf("error feeding data to session: %v", err)
		return nil, err
	}

	if offset < c.BufferLength() {
		c.ShiftN(offset)
	} else {
		c.ResetBuffer()
	}

	return nil, nil
}
