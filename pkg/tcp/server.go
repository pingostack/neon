package tcp

import (
	"github.com/panjf2000/gnet"
	"github.com/sirupsen/logrus"
)

type IManager interface {
	NewOrGet() IContext
	OnTcpClose(ctx IContext) error
}

type IContext interface {
	SetConn(c gnet.Conn)
	OnTcpClose() error
	OnTcpRread(buf []byte) (int, error)
}

type Context struct {
	Conn gnet.Conn
}

type Server struct {
	*gnet.EventServer
	addr    string
	manager IManager
}

func NewServer(addr string, manager IManager, opt gnet.Option) (*Server, error) {
	s := &Server{
		addr:    addr,
		manager: manager,
	}

	err := gnet.Serve(s, addr, opt)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// OnShutdown fires when the engine is being shut down, it is called right after
// all event-loops and connections are closed.
func (s *Server) OnShutdown(gs gnet.Server) {
	logrus.Errorf("Shutting down")
}

// OnOpen fires when a new connection has been opened.
// The parameter out is the return value which is going to be sent back to the peer.
func (s *Server) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	ctx := s.manager.NewOrGet()
	c.SetContext(ctx)
	ctx.SetConn(c)

	return
}

// OnClose fires when a connection has been closed.
// The parameter err is the last known connection error.
func (s *Server) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	ctx := c.Context().(IContext)
	ctx.OnTcpClose()
	s.manager.OnTcpClose(ctx)

	return
}

// OnTraffic fires when a local socket receives data from the peer.
func (s *Server) React(packet []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	//	ctx := c.Context().(IContext)

	logrus.Debugf("Received data from %s: %s", c.RemoteAddr(), string(packet))
	// buf := packet //c.Read()

	// var len int
	// var err error
	// if len, err = ctx.OnTcpRread(buf); err != nil {
	// 	logrus.Errorf("rtsp s onTraffic error, parse error: %s", err.Error())
	// 	return nil, gnet.Close
	// }

	// if len > 0 {
	// 	c.ShiftN(len)
	// }

	return nil, gnet.None
}

func (ctx *Context) Write(data []byte) error {
	return ctx.Conn.AsyncWrite(data)
}

func (ctx *Context) SetConn(c gnet.Conn) {
	ctx.Conn = c
}
