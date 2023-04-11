package gortc

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/let-light/neon/pkg/forwarder"
	"github.com/pion/ice/v2"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	RolePublisher  = 0
	RoleSubscriber = 1
)

type WebRTC struct {
	wg       *sync.WaitGroup
	ctx      context.Context
	settings WebRTCConfig
	w        WebRTCTransportConfig
}

var webrtcModule *WebRTC

func init() {
	webrtcModule = &WebRTC{}
}

func WebRTCModule() *WebRTC {
	return webrtcModule
}

func (rtc *WebRTC) NewPublisher(id string, logger *logrus.Entry) (forwarder.IPublisher, error) {
	return NewPublisher(id, rtc.w, logger)
}

func (rtc *WebRTC) NewSubscriber(id string, logger *logrus.Entry) (forwarder.ISubscriber, error) {
	return NewSubscriber(id, rtc.w, logger)
}

func (rtc *WebRTC) InitModule(ctx context.Context, wg *sync.WaitGroup) (interface{}, error) {
	rtc.wg = wg
	rtc.ctx = ctx
	return &rtc.settings, nil
}

func (rtc *WebRTC) InitCommand() ([]*cobra.Command, error) {

	return nil, nil
}

func (rtc *WebRTC) ConfigChanged() {
}

func (rtc *WebRTC) RootCommand(cmd *cobra.Command, args []string) {
	se := webrtc.SettingEngine{}
	se.DisableMediaEngineCopy(true)

	wc := rtc.settings

	if wc.ICESinglePort != 0 {
		udpListener, err := net.ListenUDP("udp", &net.UDPAddr{
			IP:   net.IP{0, 0, 0, 0},
			Port: wc.ICESinglePort,
		})
		if err != nil {
			panic(err)
		}
		se.SetICEUDPMux(webrtc.NewICEUDPMux(nil, udpListener))
	} else {
		var icePortStart, icePortEnd uint16

		if len(wc.ICEPortRange) == 2 {
			icePortStart = wc.ICEPortRange[0]
			icePortEnd = wc.ICEPortRange[1]
		}

		if icePortStart != 0 || icePortEnd != 0 {
			if err := se.SetEphemeralUDPPortRange(icePortStart, icePortEnd); err != nil {
				panic(err)
			}
		}
	}

	var iceServers []webrtc.ICEServer
	if wc.Candidates.IceLite {
		se.SetLite(wc.Candidates.IceLite)
	} else {
		for _, iceServer := range wc.ICEServers {
			s := webrtc.ICEServer{
				URLs:       iceServer.URLs,
				Username:   iceServer.Username,
				Credential: iceServer.Credential,
			}
			iceServers = append(iceServers, s)
		}
	}

	//se.BufferFactory = wc.BufferFactory.GetOrNew

	sdpSemantics := webrtc.SDPSemanticsUnifiedPlan
	switch wc.SDPSemantics {
	case "unified-plan-with-fallback":
		sdpSemantics = webrtc.SDPSemanticsUnifiedPlanWithFallback
	case "plan-b":
		sdpSemantics = webrtc.SDPSemanticsPlanB
	}

	if wc.Timeouts.ICEDisconnectedTimeout == 0 &&
		wc.Timeouts.ICEFailedTimeout == 0 &&
		wc.Timeouts.ICEKeepaliveInterval == 0 {
	} else {
		se.SetICETimeouts(
			time.Duration(wc.Timeouts.ICEDisconnectedTimeout)*time.Second,
			time.Duration(wc.Timeouts.ICEFailedTimeout)*time.Second,
			time.Duration(wc.Timeouts.ICEKeepaliveInterval)*time.Second,
		)
	}

	w := WebRTCTransportConfig{
		Configuration: webrtc.Configuration{
			ICEServers:   iceServers,
			SDPSemantics: sdpSemantics,
		},
		Setting: se,
		//BufferFactory: wc.BufferFactory,
	}

	if len(wc.Candidates.NAT1To1IPs) > 0 {
		w.Setting.SetNAT1To1IPs(wc.Candidates.NAT1To1IPs, webrtc.ICECandidateTypeHost)
	}

	if !wc.MDNS {
		w.Setting.SetICEMulticastDNSMode(ice.MulticastDNSModeDisabled)
	}

	rtc.w = w
}

func (rtc *WebRTC) WebRTCConfig() WebRTCTransportConfig {
	return rtc.w
}
