package rtsp

import (
	"strconv"
	"time"

	goPool "github.com/panjf2000/gnet/pkg/pool/goroutine"
	"github.com/pingopenstack/neon/src/core"
	"github.com/sirupsen/logrus"
)

const (
	EmptyState State = iota
	OptionsState
	DescribeState
	SetupState
	PlayState
	PauseState
	TeardownState
)

type Serv struct {
	logger      *logrus.Entry
	state       State
	cseqCounter int
	pool        *goPool.Pool
	write       WriteHandler
}

type IServImpl interface {
}

func NewServ(logger *logrus.Entry, write WriteHandler) *Serv {

	return &Serv{
		state:       EmptyState,
		cseqCounter: 0,
		pool:        goPool.Default(),
		write:       write,
		logger:      logger.WithField("Role", "serv"),
	}
}

func (serv *Serv) State() State {
	return serv.state
}

func (serv *Serv) decodeRtpRtcp(buf []byte) (*core.AVFrame, int, error) {
	return nil, 0, nil
}

func (serv *Serv) Feed(buf []byte) (int, error) {

	if len(buf) == 0 {
		serv.logger.Warningf("rtsp feed empty data")
		return 0, nil
	}

	var err error
	var endOffset int
	if buf[0] != '$' {
		var req *Request
		req, endOffset, err = UnmarshalRequest(buf)
		if err != nil {
			return endOffset, err
		}

		err := serv.handleRequest(req)
		if err != nil {
			return endOffset, err
		}

	} else {
		_, endOffset, err = serv.decodeRtpRtcp(buf)
		if err != nil {
			return endOffset, err
		}

	}

	return endOffset, nil
}

func (serv *Serv) handleRequest(req *Request) error {
	serv.pool.Submit(func() {
		serv.logger.Debugf("rtsp request: %s", req.String())

		var err error
		var resp IResponse
		switch req.Method() {
		case OptionsMethod:
			resp, err = serv.OptionsProcess(req)
		case DescribeMethod:

		case AnnounceMethod:
		case SetupMethod:
		case PlayMethod:
		case PauseMethod:
		case TeardownMethod:
		case GetParameterMethod:
		case SetParameterMethod:
		default:
		}

		if err != nil {
			serv.logger.Errorf("rtsp request error: %s", err.Error())
			return
		}

		if resp != nil {
			err = serv.write([]byte(resp.String()))
			if err != nil {
				serv.logger.Errorf("write rtsp response failed: %s", err.Error())
			}
		} else {
			serv.logger.Errorf("resp is nil")
		}
	})
	return nil
}

func (serv *Serv) OptionsProcess(req *Request) (IResponse, error) {
	return serv.NewOptionsResponse(req.CSeq(), []string{
		"OPTIONS",
		"ANNOUNCE",
		"DESCRIBE",
		"SETUP",
		"TEARDOWN",
		"PLAY",
		"PAUSE",
		"GET_PARAMETER",
		"SET_PARAMETER",
	}), nil
}

func (serv *Serv) NewResponse(cseq int, status Status) *Response {
	resp := &Response{
		version:   "RTSP/1.0",
		status:    status,
		statusStr: status.String(),
		lines: HeaderLines{
			"CSeq":           strconv.Itoa(cseq),
			"Date":           time.Now().Format(time.RFC1123),
			"Content-Length": strconv.Itoa(0),
			"Server":         "Neon-RTSP",
		},
	}

	return resp
}

func (serv *Serv) NewOptionsResponse(cseq int, options []string) *OptionsResponse {
	resp := serv.NewResponse(cseq, StatusOK).Option()
	resp.SetOptions(options)
	return resp
}
