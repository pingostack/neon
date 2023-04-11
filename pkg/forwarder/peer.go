package forwarder

import (
	"fmt"
	"sync"

	"github.com/let-light/neon/pkg/utils"
	"github.com/lucsky/cuid"
	"github.com/sirupsen/logrus"
)

// type ITrack interface {
// 	MarshalTracks() string
// 	SetRemoteDescription(desc string) (string, error)
// }

// type ITransport interface {
// 	MarshalTransport() string
// 	SetRemoteTransport(data string) error
// }

type IPublisher interface {
	//	OnConnect(f func(data string, publisher IPublisher))
	Close() error
	// Marshal() string
	// OnRemoteTrack(f func())
	// ITrack
	// ITransport
}

type ISubscriber interface {
	//	OnConnect(f func(data string, subscriber ISubscriber))
	Close() error
	// Marshal() string
	// OnLocalTrack(f func())
	// ITrack
	// ITransport
}

type IPeer interface {
	ID() string
	Group() IGroup
	Publisher() IPublisher
	Subscriber() ISubscriber
	Close() error
}

type IGroupProvider interface {
	GetOrNewGroup(groupId string) IGroup
}

type IPublisherProvider interface {
	NewPublisher(id string, logger *logrus.Entry) (IPublisher, error)
}

type ISubscriberProvider interface {
	NewSubscriber(id string, logger *logrus.Entry) (ISubscriber, error)
}

type JoinConfig struct {
	NoPublish        bool
	NoSubscribe      bool
	NoAutioSubscribe bool
}

type Peer struct {
	sync.Mutex
	closed             utils.AtomicBool
	id                 string
	publisherProvider  IPublisherProvider
	subscriberProvider ISubscriberProvider
	groupProvider      IGroupProvider
	publisher          IPublisher
	subscriber         ISubscriber
	group              IGroup
	logger             *logrus.Entry
	OnConnect          func()
	OnRemoteTrack      func()
	OnLocalTrack       func()
}

func NewPeer(id string, logger *logrus.Entry,
	publisherProvider IPublisherProvider,
	subscriberProvider ISubscriberProvider,
	groupProvider IGroupProvider) *Peer {

	if id == "" {
		id = cuid.New()
	}

	p := &Peer{
		publisherProvider:  publisherProvider,
		subscriberProvider: subscriberProvider,
		groupProvider:      groupProvider,
		logger: logger.WithFields(logrus.Fields{
			"peer":   id,
			"object": "peer",
		}),
	}

	p.Logger().Info("new peer created")

	return p
}

func (p *Peer) Logger() *logrus.Entry {
	return p.logger
}

func (p *Peer) Join(groupId string, config JoinConfig) error {
	p.Logger().Info("joining group")

	if p.closed.Get() {
		return fmt.Errorf("peer is closed")
	}

	p.group = p.groupProvider.GetOrNewGroup(groupId)
	if p.group == nil {
		return fmt.Errorf("failed to get group")
	}

	if p.publisherProvider != nil && !config.NoPublish {
		publisher, err := p.publisherProvider.NewPublisher(p.id, p.Logger().WithField("target", "publisher"))

		if err != nil {
			p.Logger().WithError(err).Error("failed to create publisher")
			return err
		}

		p.publisher = publisher
	}

	if p.subscriberProvider != nil && !config.NoSubscribe {
		subscriber, err := p.subscriberProvider.NewSubscriber(p.id, p.Logger().WithField("target", "subscriber"))

		if err != nil {
			p.Logger().WithError(err).Error("failed to create subscriber")
			return err
		}

		p.subscriber = subscriber
	}

	p.group.AddPeer(p)

	if !config.NoSubscribe {
		p.group.Subscribe(p)
	}

	return nil
}

func (p *Peer) ID() string {
	return p.id
}

func (p *Peer) Group() IGroup {
	return p.group
}

func (p *Peer) Publisher() IPublisher {
	return p.publisher
}

func (p *Peer) Subscriber() ISubscriber {
	return p.subscriber
}

func (p *Peer) Close() error {
	p.Lock()
	defer p.Unlock()

	if !p.closed.Set(true) {
		return nil
	}

	if p.group != nil {
		p.group.RemovePeer(p)
	}

	if p.publisher != nil {
		p.publisher.Close()
	}
	if p.subscriber != nil {
		if err := p.subscriber.Close(); err != nil {
			return err
		}
	}
	return nil
}
