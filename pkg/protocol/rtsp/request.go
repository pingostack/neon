package rtsp

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type MethodEnum int

const (
	UnknownMethod MethodEnum = iota - 1
	OptionsMethod
	DescribeMethod
	AnnounceMethod
	SetupMethod
	PlayMethod
	PauseMethod
	TeardownMethod
	GetParameterMethod
	SetParameterMethod
	RecordMethod
)

type Request struct {
	method  string
	url     string
	version string
	lines   HeaderLines
	content []byte
}

func UnmarshalRequest(buf []byte) (*Request, int, error) {
	headerEndOffset := bytes.Index(buf, []byte("\r\n\r\n"))
	if headerEndOffset == -1 {
		return nil, -1, fmt.Errorf("incomplete packet, %s", string(buf))
	}

	headerEndOffset += 2
	endOffset := headerEndOffset + 2

	contentLength := 0

	req := &Request{
		lines: make(HeaderLines),
	}

	lines := bytes.Split(buf[:headerEndOffset], []byte("\r\n"))
	if len(lines) < 2 {
		return nil, endOffset, errors.New("read lines error, invalid packet")
	}

	// parse first line
	methodLine := lines[0]
	methodLineParts := bytes.Split(methodLine, []byte(" "))
	if len(methodLineParts) != 3 {

		return nil, endOffset, fmt.Errorf("invalid method line: %s", string(methodLine))
	}

	req.method = strings.ToLower(string(methodLineParts[0]))
	req.url = strings.ToLower(string(methodLineParts[1]))
	req.version = strings.ToLower(string(methodLineParts[2]))

	// parse other lines
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		idx := bytes.Index(line, []byte(":"))
		if idx == -1 {
			continue
		}

		key := strings.ToLower(string(line[:idx]))
		value := string(line[idx+1:])
		value = strings.TrimSpace(value)
		req.lines[key] = value

		if key == "content-length" {
			var err error
			contentLength, err = strconv.Atoi(value)
			if err != nil {
				return nil, headerEndOffset + 4, err
			}
		}
	}

	req.content = buf[headerEndOffset+4 : headerEndOffset+4+contentLength]

	endOffset += contentLength

	return req, endOffset, nil
}

func (lines HeaderLines) String() string {
	s := ""
	for k, v := range lines {
		s += k + ": " + v + "\r\n"
	}
	return s
}

func (req *Request) Method() MethodEnum {
	switch req.method {
	case "options":
		return OptionsMethod
	case "describe":
		return DescribeMethod
	case "announce":
		return AnnounceMethod
	case "setup":
		return SetupMethod
	case "play":
		return PlayMethod
	case "pause":
		return PauseMethod
	case "teardown":
		return TeardownMethod
	case "getparameter":
		return GetParameterMethod
	case "setparameter":
		return SetParameterMethod
	case "record":
		return RecordMethod
	default:
		return UnknownMethod
	}
}

func (req *Request) MethodStr() string {
	return req.method
}

func (req *Request) GetLine(key string) string {
	return req.lines[key]
}

func (req *Request) SetLine(key, value string) {
	req.lines[key] = value
}

func (req *Request) String() string {
	return req.method + " " + req.url + " " + req.version + "\r\n" + req.lines.String() + "\r\n" + string(req.content)
}

func (req *Request) CSeq() int {
	cseqLine := req.lines["cseq"]
	if cseqLine == "" {
		return -1
	}

	cseq, err := strconv.Atoi(cseqLine)
	if err != nil {
		return -1
	}

	return cseq
}

func (req *Request) Session() string {
	return req.lines["session"]
}

func (req *Request) Option() *OptionsRequest {
	return &OptionsRequest{
		Request: req,
	}
}

func (req *Request) Describe() *DescribeRequest {
	return &DescribeRequest{
		Request: req,
	}
}

func (req *Request) Announce() *AnnounceRequest {
	return &AnnounceRequest{
		Request: req,
	}
}

func (req *Request) Setup() *SetupRequest {
	return &SetupRequest{
		Request: req,
	}
}

func (req *Request) Play() *PlayRequest {
	return &PlayRequest{
		Request: req,
	}
}

func (req *Request) Pause() *PauseRequest {
	return &PauseRequest{
		Request: req,
	}
}

func (req *Request) Teardown() *TeardownRequest {
	return &TeardownRequest{
		Request: req,
	}
}

func (req *Request) GetParameter() *GetParameterRequest {
	return &GetParameterRequest{
		Request: req,
	}
}

func (req *Request) SetParameter() *SetParameterRequest {
	return &SetParameterRequest{
		Request: req,
	}
}

func (req *Request) Record() *RecordRequest {
	return &RecordRequest{
		Request: req,
	}
}

// OptionsRequest is a RTSP OPTIONS request
type OptionsRequest struct {
	*Request
}

func (req *OptionsRequest) Require() string {
	return req.lines["require"]
}

func (req *OptionsRequest) ProxyRequire() string {
	return req.lines["proxy-require"]
}

// DescribeRequest is a RTSP DESCRIBE request
type DescribeRequest struct {
	*Request
}

func (req *DescribeRequest) Accept() string {
	return req.lines["accept"]
}

func (req *DescribeRequest) SetSdp(sdp string) {
	req.lines["content-type"] = "application/sdp"
	req.content = []byte(sdp)
}

// AnnounceRequest is a RTSP ANNOUNCE request
type AnnounceRequest struct {
	*Request
}

func (req *AnnounceRequest) ContentType() string {
	return req.lines["content-type"]
}

// SetupRequest is a RTSP SETUP request
type SetupRequest struct {
	*Request
}

func (req *SetupRequest) TransportString() string {
	return req.lines["transport"]
}

func (req *SetupRequest) Transport() (*Transport, error) {
	return UnmarshalTransport(req.TransportString())
}

func (req *SetupRequest) SetTransport(transport *Transport) {
	req.lines["transport"], _ = transport.String()
}

// PlayRequest is a RTSP PLAY request
type PlayRequest struct {
	*Request
}

func (req *PlayRequest) Range() string {
	return req.lines["range"]
}

// PauseRequest is a RTSP PAUSE request
type PauseRequest struct {
	*Request
}

func (req *PauseRequest) Range() string {
	return req.lines["range"]
}

// TeardownRequest is a RTSP TEARDOWN request
type TeardownRequest struct {
	*Request
}

// GetParameterRequest is a RTSP GET_PARAMETER request
type GetParameterRequest struct {
	*Request
}

func (req *GetParameterRequest) Parameters() []string {
	if len(req.content) == 0 {
		return nil
	}

	return strings.Split(string(req.content), "\r\n")
}

// SetParameterRequest is a RTSP SET_PARAMETER request
type SetParameterRequest struct {
	*Request
}

func (req *SetParameterRequest) Parameters() map[string]string {
	if len(req.content) == 0 {
		return nil
	}

	lines := strings.Split(string(req.content), "\r\n")
	params := make(map[string]string)
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		params[parts[0]] = parts[1]
	}

	return params
}

// RecordRequest is a RTSP RECORD request
type RecordRequest struct {
	*Request
}
