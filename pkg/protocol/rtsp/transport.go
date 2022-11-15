package rtsp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type TransportMode int
type RtpProfile int

const (
	TransportModeTcp TransportMode = iota
	TransportModeUdp
)

const (
	RtpProfileInvalid = iota - 1 // invalid
	RtpProfileAVP
	RtpProfileAVPF
	RtpProfileSAVP
	RtpProfileSAVPF
)

const (
	RtpProfileAVPStr   = "RTP/AVP"
	RtpProfileAVPFStr  = "RTP/AVPF"
	RtpProfileSAVPStr  = "RTP/SAVP"
	RtpProfileSAVPFStr = "RTP/SAVPF"
)

type Transport struct {
	Profile      RtpProfile
	Mode         TransportMode
	Interleaveds []int
	Unicast      bool
	ClientPorts  []int
	ServerPorts  []int
	SSRC         int64
}

func (r RtpProfile) ToString() string {
	switch r {
	case RtpProfileAVP:
		return RtpProfileAVPStr
	case RtpProfileAVPF:
		return RtpProfileAVPFStr
	case RtpProfileSAVP:
		return RtpProfileSAVPStr
	case RtpProfileSAVPF:
		return RtpProfileSAVPFStr
	}

	return ""
}

func (r *RtpProfile) Parse(s string) error {
	if strings.Contains(s, RtpProfileAVPStr) {
		*r = RtpProfileAVP
		return nil
	} else if strings.Contains(s, RtpProfileAVPFStr) {
		*r = RtpProfileAVPF
		return nil
	} else if strings.Contains(s, RtpProfileSAVPStr) {
		*r = RtpProfileSAVP
		return nil
	} else if strings.Contains(s, RtpProfileSAVPFStr) {
		*r = RtpProfileSAVPF
		return nil
	}

	return errors.New("invalid rtp profile")
}

func (t *Transport) String() (string, error) {
	s := "Transport: "

	if t.Profile == RtpProfileInvalid {
		return "", errors.New("invalid rtp profile")
	}

	s += t.Profile.ToString()

	if t.Mode == TransportModeTcp {
		s += "TCP;"

		s += "interleaved="
		for i, v := range t.Interleaveds {
			s += strconv.Itoa(v)
			if i < len(t.Interleaveds)-1 {
				s += "-"
			}
		}
	} else {
		s += "UDP;"
		if t.Unicast {
			s += "unicast;"
		}

		if len(t.ClientPorts) > 0 {
			s += "client_port="
			for i, v := range t.ClientPorts {
				s += strconv.Itoa(v)
				if i < len(t.ClientPorts)-1 {
					s += "-"
				}
			}
		}

		if len(t.ServerPorts) > 0 {
			s += "server_port="
			for i, v := range t.ServerPorts {
				s += strconv.Itoa(v)
				if i < len(t.ServerPorts)-1 {
					s += "-"
				}
			}
		}
	}

	if t.SSRC != 0 {
		s += fmt.Sprintf(";ssrc=%x", t.SSRC)
	}

	return s, nil
}

func NewTransport(s string) (*Transport, error) {
	t := &Transport{
		Profile:      RtpProfileInvalid,
		Mode:         TransportModeUdp,
		Interleaveds: make([]int, 0),
		Unicast:      true,
		ClientPorts:  make([]int, 0),
		ServerPorts:  make([]int, 0),
		SSRC:         0,
	}

	parts := strings.Split(s, ";")

	for _, p := range parts {
		var key, val string
		kv := strings.Split(p, "=")
		if len(kv) == 2 {
			key = strings.ToLower(kv[0])
			val = kv[1]

			switch key {
			case "interleaved":
				t.Interleaveds = make([]int, 2)
				iv := strings.Split(val, "-")
				if len(iv) == 2 {
					t.Interleaveds[0], _ = strconv.Atoi(iv[0])
					t.Interleaveds[1], _ = strconv.Atoi(iv[0])

					t.Mode = TransportModeTcp
				}

			case "client_port":
				iv := strings.Split(val, "-")
				if len(iv) == 2 {
					t.ClientPorts[0], _ = strconv.Atoi(iv[0])
					t.ClientPorts[1], _ = strconv.Atoi(iv[0])

					t.Mode = TransportModeUdp
				}

			case "server_port":
				iv := strings.Split(val, "-")
				if len(iv) == 2 {
					t.ClientPorts[0], _ = strconv.Atoi(iv[0])
					t.ClientPorts[1], _ = strconv.Atoi(iv[0])

					t.Mode = TransportModeUdp
				}

			case "ssrc":
				t.SSRC, _ = strconv.ParseInt(val, 16, 32)
			}

		} else {
			key = strings.ToLower(p)

			if strings.Contains(key, "rtp/") {
				t.Profile.Parse(key)
			} else if strings.Contains(key, "unicast") {
				t.Unicast = true
			}
		}
	}

	if t.Mode == TransportModeUdp && len(t.ClientPorts) != 2 {
		return nil, errors.New("invalid transport ports")
	}

	if t.Mode == TransportModeTcp && len(t.Interleaveds) != 2 {
		return nil, errors.New("invalid transport interleaved")
	}

	if t.Profile == RtpProfileInvalid {
		return nil, errors.New("invalid rtp profile")
	}

	if t.SSRC == 0 {
		return nil, errors.New("invalid ssrc")
	}

	return t, nil
}

func (t *Transport) RtpPort() int {
	if len(t.ClientPorts) == 0 {
		return -1
	}

	return t.ClientPorts[0]
}

func (t *Transport) RtcpPort() int {
	if len(t.ClientPorts) == 0 {
		return -1
	}

	return t.ClientPorts[1]
}

func (t *Transport) RtpInterleaved() int {
	if len(t.Interleaveds) == 0 {
		return -1
	}

	return t.Interleaveds[0]
}

func (t *Transport) RtcpInterleaved() int {
	if len(t.Interleaveds) == 0 {
		return -1
	}

	return t.Interleaveds[1]
}
