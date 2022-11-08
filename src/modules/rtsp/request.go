package rtsp

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type Request struct {
	headers map[string]string
	data    string
}

func (req *Request) Marshal() (string, error) {
	return "", nil
}

func Unmarshal(buf []byte, req *Request) (int, error) {
	headerEndOffset := bytes.IndexAny(buf, "\r\n\r\n")
	if headerEndOffset == -1 {
		return -1, errors.New("incomplete packet")
	}

	headers := make(map[string]string)
	lines := bytes.Split(buf, []byte("\r\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		idx := bytes.IndexAny(line, ":")
		if idx == -1 {
			continue
		}

		key := strings.ToLower(string(line[:idx]))
		value := string(line[idx+1:])
		headers[key] = value
	}

	var contentLength int
	var err error

	if headers["content-length"] != "" {
		contentLength, err = strconv.Atoi(headers["content-length"])
		if err != nil {
			return -1, err
		}

		if len(buf) < headerEndOffset+4+contentLength {
			return -1, errors.New("incomplete packet")
		}
	}

	return -1, nil
}
