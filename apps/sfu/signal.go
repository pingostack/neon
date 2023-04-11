package sfu

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/let-light/neon/pkg/forwarder"
	"github.com/let-light/neon/pkg/gortc"
	"github.com/pion/webrtc/v3"
	"github.com/sourcegraph/jsonrpc2"
)

// Join message sent when initializing a peer connection
type Join struct {
	GID    string                    `json:"sid"`
	UID    string                    `json:"uid"`
	Offer  webrtc.SessionDescription `json:"offer"`
	Config forwarder.JoinConfig      `json:"config"`
}

// Negotiation message sent when renegotiating the peer connection
type Negotiation struct {
	Desc webrtc.SessionDescription `json:"desc"`
}

// Trickle message sent when renegotiating the peer connection
type Trickle struct {
	Target    int                     `json:"target"`
	Candidate webrtc.ICECandidateInit `json:"candidate"`
}

type JSONSignal struct {
	*forwarder.Peer
}

func NewJSONSignal(p *forwarder.Peer) *JSONSignal {
	return &JSONSignal{p}
}

// Handle incoming RPC call events like join, answer, offer and trickle
func (p *JSONSignal) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	defer func() {
		if r := recover(); r != nil {
			p.Logger().Errorf("Recovered from panic: %v", r)
		}
	}()

	replyError := func(err error) {
		p.Logger().WithError(err).Error("error handling request")
		_ = conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    500,
			Message: fmt.Sprintf("%s", err),
		})
	}

	var publisher *gortc.Publisher
	var subscriber *gortc.Subscriber

	if p.Publisher() != nil {
		publisher = p.Publisher().(*gortc.Publisher)
	}

	if p.Subscriber() != nil {
		subscriber = p.Subscriber().(*gortc.Subscriber)
	}

	p.Logger().Debugf("Method %s, Params %s", req.Method, *req.Params)

	switch req.Method {
	case "join":
		var join Join
		err := json.Unmarshal(*req.Params, &join)
		if err != nil {
			p.Logger().Error(err, "connect: error parsing offer")
			replyError(err)
			break
		}

		err = p.Join(join.GID, join.Config)
		if err != nil {
			replyError(err)
			break
		}

		if p.Publisher() != nil {
			publisher = p.Publisher().(*gortc.Publisher)
		}

		if p.Subscriber() != nil {
			subscriber = p.Subscriber().(*gortc.Subscriber)
		}

		if subscriber != nil {
			subscriber.OnOffer = func(offer *webrtc.SessionDescription) {
				if err := conn.Notify(ctx, "offer", offer); err != nil {
					p.Logger().WithError(err).Error("error sending offer")
				}
			}

			subscriber.OnIceCandidate = func(candidate *webrtc.ICECandidateInit, target int) {
				if err := conn.Notify(ctx, "trickle", Trickle{
					Candidate: *candidate,
					Target:    target,
				}); err != nil {
					p.Logger().WithError(err).Info("error sending ice candidate")
				}
			}
		}

		if publisher != nil {
			publisher.OnIceCandidate = func(candidate *webrtc.ICECandidateInit, target int) {
				if err := conn.Notify(ctx, "trickle", Trickle{
					Candidate: *candidate,
					Target:    target,
				}); err != nil {
					p.Logger().WithError(err).Info("error sending ice candidate")
				}
			}

			answer, err := publisher.Answer(join.Offer)
			if err != nil {
				p.Logger().Error(err, "error marshaling offer")
				replyError(err)
				break
			}

			_ = conn.Reply(ctx, req.ID, answer)
		} else {
			replyError(fmt.Errorf("no publisher, cannot answer"))
		}

	case "offer":
		if publisher == nil {
			replyError(fmt.Errorf("no publisher, cannot answer"))
			break
		}

		var negotiation Negotiation
		err := json.Unmarshal(*req.Params, &negotiation)
		if err != nil {
			p.Logger().WithError(err).Error("connect: error parsing offer")
			replyError(err)
			break
		}

		answer, err := publisher.Answer(negotiation.Desc)
		if err != nil {
			p.Logger().WithError(err).Error("error marshaling offer")
			replyError(err)
			break
		}

		_ = conn.Reply(ctx, req.ID, answer)

	case "answer":
		if subscriber == nil {
			replyError(fmt.Errorf("no subscriber"))
			break
		}

		var negotiation Negotiation
		err := json.Unmarshal(*req.Params, &negotiation)
		if err != nil {
			p.Logger().WithError(err).Error("connect: error parsing answer")
			replyError(err)
			break
		}

		err = subscriber.SetRemoteDescription(negotiation.Desc)
		if err != nil {
			replyError(err)
		}

	case "trickle":
		var trickle Trickle
		err := json.Unmarshal(*req.Params, &trickle)
		if err != nil {
			p.Logger().WithError(err).Error("connect: error parsing candidate")
			replyError(err)
			break
		}

		if trickle.Target == gortc.RolePublisher {
			if publisher == nil {
				replyError(fmt.Errorf("no publisher"))
				break
			}
			err = publisher.AddICECandidate(trickle.Candidate)
		} else if trickle.Target == gortc.RoleSubscriber {
			if subscriber == nil {
				replyError(fmt.Errorf("no subscriber"))
				break
			}
			err = subscriber.AddICECandidate(trickle.Candidate)
		} else {
			err = fmt.Errorf("unknown target")
		}

		if err != nil {
			p.Logger().WithError(err).Error("error setting remote transport")
			replyError(err)
		}

	default:
		p.Logger().Error(nil, "connect: unknown method", "method", req.Method)
	}
}
