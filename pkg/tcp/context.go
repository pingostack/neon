package tcp

import (
	"github.com/panjf2000/gnet"
)

type IContext interface {
	SetConn(c interface{})
	OnTcpClose() error
	OnTcpRread(buf []byte) (int, error)
	RemoteAddr() string
	LocalAddr() string
	Network() string
}

type Context struct {
	Conn gnet.Conn
}

func NewContext(c interface{}) IContext {
	return &Context{
		Conn: c.(gnet.Conn),
	}
}

func (ctx *Context) Write(data []byte) error {
	return ctx.Conn.AsyncWrite(data)
}

func (ctx *Context) SetConn(c interface{}) {
	if ctx.Conn != nil {
		ctx.Conn = c.(gnet.Conn)
	} else {
		ctx.Conn = nil
	}
}

func (ctx *Context) OnTcpClose() error {
	return nil
}

func (ctx *Context) OnTcpRread(buf []byte) (int, error) {
	return 0, nil
}

func (ctx *Context) RemoteAddr() string {
	if (ctx.Conn == nil) || (ctx.Conn.RemoteAddr() == nil) {
		return ""
	}

	return ctx.Conn.RemoteAddr().String()
}

func (ctx *Context) LocalAddr() string {
	if (ctx.Conn == nil) || (ctx.Conn.RemoteAddr() == nil) {
		return ""
	}

	return ctx.Conn.LocalAddr().String()
}

func (ctx *Context) Network() string {
	if ctx.Conn == nil {
		return ""
	}

	return ctx.Conn.RemoteAddr().Network()
}
