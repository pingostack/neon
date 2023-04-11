package gortc

import (
	"fmt"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

type Publisher struct {
	OnIceCandidate             func(*webrtc.ICECandidateInit, int)
	OnICEConnectionStateChange func(webrtc.ICEConnectionState)

	id         string
	logger     *logrus.Entry
	pc         *webrtc.PeerConnection
	cfg        *WebRTCTransportConfig
	candidates []webrtc.ICECandidateInit
}

func NewPublisher(id string, cfg WebRTCTransportConfig, logger *logrus.Entry) (*Publisher, error) {
	me, err := getPublisherMediaEngine()
	if err != nil {
		logger.WithError(err).Error("NewPublisher error, getPublisherMediaEngine")
		return nil, err
	}

	logger.Infof("NewPublisher, cfg %+v", cfg.Configuration)

	api := webrtc.NewAPI(webrtc.WithMediaEngine(me), webrtc.WithSettingEngine(cfg.Setting))
	pc, err := api.NewPeerConnection(cfg.Configuration)
	if err != nil {
		logger.WithError(err).Error("NewPublisher error, NewPeerConnection")
		return nil, err
	}

	publisher := &Publisher{
		id:     id,
		logger: logger,
		pc:     pc,
		cfg:    &cfg,
	}

	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		logger.Info("on publisher ice candidate called for peer")
		if c == nil {
			return
		}

		if publisher.OnIceCandidate != nil {
			json := c.ToJSON()
			publisher.OnIceCandidate(&json, RolePublisher)
		}
	})

	pc.OnICEConnectionStateChange(func(s webrtc.ICEConnectionState) {
		logger.Info("on publisher ice connection state change called for peer ", "state ", s)
	})

	return publisher, nil
}

func (p *Publisher) AddICECandidate(candidate webrtc.ICECandidateInit) error {
	p.Logger().Info("add publisher ice candidate called for peer")
	if p.pc.RemoteDescription() != nil {
		p.Logger().Info("add publisher ice candidate called for peer, remote description is not nil")
		return p.pc.AddICECandidate(candidate)
	}
	p.candidates = append(p.candidates, candidate)
	return nil
}

func (p *Publisher) Answer(sdp webrtc.SessionDescription) (*webrtc.SessionDescription, error) {

	p.Logger().Infof("got offer sdp: %s", sdp.SDP)

	if p.pc.SignalingState() != webrtc.SignalingStateStable {
		return nil, fmt.Errorf("signalstate is not stable")
	}

	if err := p.pc.SetRemoteDescription(sdp); err != nil {
		return nil, err
	}

	for _, c := range p.candidates {
		if err := p.pc.AddICECandidate(c); err != nil {
			p.Logger().WithError(err).Error("Add publisher ice candidate to peer err")
		}
	}

	p.Logger().Infof("publisher candidates: %v", p.candidates)
	p.candidates = nil

	answer, err := p.pc.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	if err := p.pc.SetLocalDescription(answer); err != nil {
		return nil, err
	}

	return &answer, nil
}

func (p *Publisher) Logger() *logrus.Entry {
	return p.logger
}

// IPublisher impl begin

func (p *Publisher) Close() error {
	return nil
}

// IPublisher impl end
