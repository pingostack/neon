package sfu

import (
	"context"
	"net/http"

	"github.com/gogf/gf/util/guid"
	"github.com/gorilla/websocket"
	"github.com/let-light/neon/pkg/forwarder"
	"github.com/let-light/neon/pkg/gortc"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
	websocketjsonrpc2 "github.com/sourcegraph/jsonrpc2/websocket"
)

type SignalServerSettings struct {
	listenAddr string
	cert       string
	key        string
}

type SignalServer struct {
	settings SignalServerSettings
	onClose  func()
	ctx      context.Context
}

func NewSignalServer(ctx context.Context, settings SignalServerSettings) *SignalServer {
	return &SignalServer{
		settings: settings,
		ctx:      ctx,
	}
}

func (ss *SignalServer) Run() error {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.Handle("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		defer c.Close()

		peerId := guid.S()
		p := NewJSONSignal(
			forwarder.NewPeer(guid.S(),
				logrus.WithFields(logrus.Fields{
					"peer":   peerId,
					"entry":  "sfu",
					"remote": r.RemoteAddr,
					"local":  r.Host,
				}),
				gortc.WebRTCModule(),
				gortc.WebRTCModule(),
				forwarder.ForwardingSys,
			))
		defer p.Close()

		jc := jsonrpc2.NewConn(r.Context(), websocketjsonrpc2.NewObjectStream(c), p)
		<-jc.DisconnectNotify()
	}))

	var err error
	addr := ss.settings.listenAddr
	cert := ss.settings.cert
	key := ss.settings.key

	if cert != "" && key != "" {
		logrus.Info("Started listening addr https://" + addr)
		err = http.ListenAndServeTLS(addr, cert, key, nil)
	} else {
		logrus.Info("Started listening addr http://" + addr)
		err = http.ListenAndServe(addr, nil)
	}
	if err != nil {
		logrus.Error("Error starting signalserver", err)
	}

	return nil
}

func (ss *SignalServer) Shutdown() error {
	return nil
}

func (ss *SignalServer) OnClose(f func()) {
	ss.onClose = f
}
