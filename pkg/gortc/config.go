package gortc

import (
	"github.com/pion/webrtc/v3"
)

type Candidates struct {
	IceLite    bool     `mapstructure:"icelite"`
	NAT1To1IPs []string `mapstructure:"nat1to1"`
}

// WebRTCTransportConfig represents Configuration options
type WebRTCTransportConfig struct {
	Configuration   webrtc.Configuration
	Setting         webrtc.SettingEngine
	MaxBitRate      int
	TrackingPackets int
	//	BufferFactory *buffer.Factory
}

type WebRTCTimeoutsConfig struct {
	ICEDisconnectedTimeout int `mapstructure:"disconnected"`
	ICEFailedTimeout       int `mapstructure:"failed"`
	ICEKeepaliveInterval   int `mapstructure:"keepalive"`
}

// ICEServerConfig defines parameters for ice servers
type ICEServerConfig struct {
	URLs       []string `mapstructure:"urls"`
	Username   string   `mapstructure:"username"`
	Credential string   `mapstructure:"credential"`
}

// WebRTCConfig defines parameters for ice
type WebRTCConfig struct {
	ICESinglePort   int                  `mapstructure:"singleport"`
	ICEPortRange    []uint16             `mapstructure:"portrange"`
	ICEServers      []ICEServerConfig    `mapstructure:"iceserver"`
	Candidates      Candidates           `mapstructure:"candidates"`
	SDPSemantics    string               `mapstructure:"sdpsemantics"`
	MDNS            bool                 `mapstructure:"mdns"`
	Timeouts        WebRTCTimeoutsConfig `mapstructure:"timeouts"`
	TrackingPackets int                  `mapstructure:"trackingPackets"`
	MaxBitRate      int                  `mapstructure:"maxBitRate"`
	//	BufferFactory   *buffer.Factory
}

func (wc *WebRTCConfig) Compare(other *WebRTCConfig) bool {
	if wc.ICESinglePort != other.ICESinglePort {
		return false
	}

	if len(wc.ICEPortRange) != len(other.ICEPortRange) {
		return false
	}

	for i := range wc.ICEPortRange {
		if wc.ICEPortRange[i] != other.ICEPortRange[i] {
			return false
		}
	}

	if len(wc.ICEServers) != len(other.ICEServers) {
		return false
	}

	for i := range wc.ICEServers {
		if wc.ICEServers[i].URLs[0] != other.ICEServers[i].URLs[0] {
			return false
		}
		if wc.ICEServers[i].Username != other.ICEServers[i].Username {
			return false
		}
		if wc.ICEServers[i].Credential != other.ICEServers[i].Credential {
			return false
		}
	}

	if wc.Candidates.IceLite != other.Candidates.IceLite {
		return false
	}

	if len(wc.Candidates.NAT1To1IPs) != len(other.Candidates.NAT1To1IPs) {
		return false
	}

	for i := range wc.Candidates.NAT1To1IPs {
		if wc.Candidates.NAT1To1IPs[i] != other.Candidates.NAT1To1IPs[i] {
			return false
		}
	}

	if wc.SDPSemantics != other.SDPSemantics {
		return false
	}

	if wc.MDNS != other.MDNS {
		return false
	}

	if wc.Timeouts.ICEDisconnectedTimeout != other.Timeouts.ICEDisconnectedTimeout {
		return false
	}

	if wc.Timeouts.ICEFailedTimeout != other.Timeouts.ICEFailedTimeout {
		return false
	}

	if wc.Timeouts.ICEKeepaliveInterval != other.Timeouts.ICEKeepaliveInterval {
		return false
	}

	if wc.MaxBitRate != other.MaxBitRate {
		return false
	}

	return true
}
