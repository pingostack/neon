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
	logger    *logrus.Entry
}

func NewRtspServer(settings RtspServerSettings, logger *logrus.Entry) (*RtspServer, error) {
	var err error
	server := &RtspServer{
		logger: logger,
	}

	server.settings = settings

	settings.TcpSettings.Logger = logrus.StandardLogger()
	settings.TcpSettings.Codec = server

	server.tcpServer, err = tcp.NewServer(settings.Addr, server, gnet.WithOptions(settings.TcpSettings))
	if err != nil {
		Logger().Errorf("error creating tcp server: %v", err)
		return nil, err
	}

	return server, nil
}

func (server *RtspServer) NewOrGet() tcp.IContext {
	session := NewSession(server.logger)

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
	ctx := c.Context().(tcp.IContext)
	buf := c.Read()
	var len int
	var err error
	if len, err = ctx.OnTcpRread(buf); err != nil {
		logrus.Errorf("rtsp s onTraffic error, parse error: %s", err.Error())
		return nil, nil
	}

	if len > 0 {
		logrus.Debugf("rtsp s onTraffic, len: %d", len)
		// 	c.ShiftN(len)
	}

	return nil, nil
}
