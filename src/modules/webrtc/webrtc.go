package webrtc

import (
	"errors"
	"net"
	"runtime"
	"time"

	"github.com/let-light/neon/pkg/module"
	"github.com/pion/ice/v2"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ICEServerConfig struct {
	URLs       []string `mapstructure:"urls"`
	Username   string   `mapstructure:"username"`
	Credential string   `mapstructure:"credential"`
}

type ICEConfig struct {
	ICESinglePort          int               `mapstructure:"singleport"`
	ICEPortRange           []uint16          `mapstructure:"portrange"`
	ICEServers             []ICEServerConfig `mapstructure:"iceserver"`
	IceLite                bool              `mapstructure:"icelite"`
	NAT1To1IPs             []string          `mapstructure:"nat1to1"`
	SDPSemantics           string            `mapstructure:"sdpsemantics"`
	MDNS                   bool              `mapstructure:"mdns"`
	ICEDisconnectedTimeout int               `mapstructure:"disconnectedtimeout"`
	ICEFailedTimeout       int               `mapstructure:"failedtimeout"`
	ICEKeepaliveInterval   int               `mapstructure:"keepaliveinterval"`
}

type AudioLevelObserverConfig struct {
	AudioLevelInterval  int   `mapstructure:"audiolevelinterval"`
	AudioLevelThreshold uint8 `mapstructure:"audiolevelthreshold"`
	AudioLevelFilter    int   `mapstructure:"audiolevelfilter"`
}

type TrackConfig struct {
	MaxBandwidth   uint64 `mapstructure:"maxbandwidth"`
	MaxPacketTrack int    `mapstructure:"maxpackettrack"`
}

type WebRTCFlags struct {
}

type WebRTCSettings struct {
	Ballast             int64                    `mapstructure:"ballast"`
	WithStats           bool                     `mapstructure:"withstats"`
	ICE                 ICEConfig                `mapstructure:"ice"`
	AudioLevelObserver  AudioLevelObserverConfig `mapstructure:"audiolevelobserver"`
	Track               TrackConfig              `mapstructure:"track"`
	WebRTCConfiguration webrtc.Configuration
	WebRTCSettingEngine webrtc.SettingEngine
}

type WebRTCModule struct {
	flags    *WebRTCFlags
	settings *WebRTCSettings
}

var instance *WebRTCModule

func init() {
	instance = &WebRTCModule{
		flags:    &WebRTCFlags{},
		settings: &WebRTCSettings{},
	}
	module.Register(instance)
}

func WebRTCModuleInstance() module.IModule {
	return instance
}

func (w *WebRTCModule) OnInitModule() (interface{}, error) {
	return w.settings, nil
}

func (w *WebRTCModule) OnInitCommand() ([]*cobra.Command, error) {
	return nil, nil
}

func (w *WebRTCModule) OnConfigModified() {

}

func (w *WebRTCModule) OnPostInitCommand() {

}

func (w *WebRTCModule) OnMainRun(cmd *cobra.Command, args []string) {
	w.initWebRTCEngineConfig()
}

func (w *WebRTCModule) initWebRTCEngineConfig() {
	if w.settings.Ballast != 0 {
		ballast := make([]byte, w.settings.Ballast*1024*1024)
		runtime.KeepAlive(ballast)
		logrus.Info("Ballast memory size: ", w.settings.Ballast*1024*1024, "(bytes)")
	}

	se := webrtc.SettingEngine{}
	se.DisableMediaEngineCopy(true)

	if w.settings.ICE.ICESinglePort != 0 {
		logrus.Info("Listen on single-port: ", w.settings.ICE.ICESinglePort)
		udpListener, err := net.ListenUDP("udp", &net.UDPAddr{
			IP:   net.IP{0, 0, 0, 0},
			Port: w.settings.ICE.ICESinglePort,
		})
		if err != nil {
			panic(err)
		}
		se.SetICEUDPMux(webrtc.NewICEUDPMux(nil, udpListener))
	} else {
		var icePortStart, icePortEnd uint16

		if len(w.settings.ICE.ICEPortRange) == 2 {
			icePortStart = w.settings.ICE.ICEPortRange[0]
			icePortEnd = w.settings.ICE.ICEPortRange[1]
		}

		if icePortStart != 0 || icePortEnd != 0 {
			if err := se.SetEphemeralUDPPortRange(icePortStart, icePortEnd); err != nil {
				panic(err)
			}
		} else {
			panic(errors.New("missing ice port"))
		}
	}

	var iceServers []webrtc.ICEServer
	if w.settings.ICE.IceLite {
		se.SetLite(w.settings.ICE.IceLite)
	} else {
		for _, iceServer := range w.settings.ICE.ICEServers {
			s := webrtc.ICEServer{
				URLs:       iceServer.URLs,
				Username:   iceServer.Username,
				Credential: iceServer.Credential,
			}
			iceServers = append(iceServers, s)
		}
	}

	se.BufferFactory = nil

	sdpSemantics := webrtc.SDPSemanticsUnifiedPlan
	switch w.settings.ICE.SDPSemantics {
	case "unified-plan-with-fallback":
		sdpSemantics = webrtc.SDPSemanticsUnifiedPlanWithFallback
	case "plan-b":
		sdpSemantics = webrtc.SDPSemanticsPlanB
	}

	if w.settings.ICE.ICEDisconnectedTimeout == 0 &&
		w.settings.ICE.ICEFailedTimeout == 0 &&
		w.settings.ICE.ICEKeepaliveInterval == 0 {
		logrus.Info("No webrtc timeouts found in config, using default ones")
	} else {
		se.SetICETimeouts(
			time.Duration(w.settings.ICE.ICEDisconnectedTimeout)*time.Second,
			time.Duration(w.settings.ICE.ICEFailedTimeout)*time.Second,
			time.Duration(w.settings.ICE.ICEKeepaliveInterval)*time.Second,
		)
	}
	w.settings.WebRTCSettingEngine = se
	w.settings.WebRTCConfiguration = webrtc.Configuration{
		ICEServers:   iceServers,
		SDPSemantics: sdpSemantics,
	}

	if len(w.settings.ICE.NAT1To1IPs) > 0 {
		w.settings.WebRTCSettingEngine.SetNAT1To1IPs(w.settings.ICE.NAT1To1IPs, webrtc.ICECandidateTypeHost)
	}

	if !w.settings.ICE.MDNS {
		w.settings.WebRTCSettingEngine.SetICEMulticastDNSMode(ice.MulticastDNSModeDisabled)
	}
}

func (w *WebRTCModule) GetWebRTCConfiguration() *webrtc.Configuration {
	return &w.settings.WebRTCConfiguration
}

func (w *WebRTCModule) GetWebRTCSettingsEngine() *webrtc.SettingEngine {
	return &w.settings.WebRTCSettingEngine
}
