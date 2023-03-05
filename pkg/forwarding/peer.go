package forwarding

import (
	"fmt"
	"sync"

	"github.com/let-light/neon/pkg/utils"
	"github.com/lucsky/cuid"
	"github.com/sirupsen/logrus"
)

type IPublisher interface {
	Close() error
}

type ISubscriber interface {
	Close() error
}

type IPeer interface {
	ID() string
	RouterGroup() IRouterGroup
	Publisher() IPublisher
	Subscriber() ISubscriber
	Close() error
}

type IRouterGroupProvider interface {
	GetOrNewRouterGroup(routerGroupId string) IRouterGroup
}

type IPublisherProvider interface {
	NewPublisher(id string, routerGroup IRouterGroup) IPublisher
}

type ISubscriberProvider interface {
	NewSubscriber(id string, routerGroup IRouterGroup) ISubscriber
}

type JoinConfig struct {
	NoPublish        bool
	NoSubscribe      bool
	NoAutioSubscribe bool
}

type Peer struct {
	sync.Mutex
	closed              utils.AtomicBool
	id                  string
	publisherProvider   IPublisherProvider
	subscriberProvider  ISubscriberProvider
	routerGroupProvider IRouterGroupProvider
	publisher           IPublisher
	subscriber          ISubscriber
	routerGroup         IRouterGroup
	logger              *logrus.Entry
	OnTrack             func()
}

func NewPeer(id string, logger *logrus.Entry,
	publisherProvider IPublisherProvider,
	subscriberProvider ISubscriberProvider,
	routerGroupProvider IRouterGroupProvider) *Peer {

	if id == "" {
		id = cuid.New()
	}

	p := &Peer{
		publisherProvider:   publisherProvider,
		subscriberProvider:  subscriberProvider,
		routerGroupProvider: routerGroupProvider,
		logger:              logger.WithFields(logrus.Fields{"peer": id, "module": "peer"}),
	}

	p.logger.Info("new peer created")

	return p
}

func (p *Peer) Join(routerGroupId string, config JoinConfig) error {
	p.Lock()
	defer p.Unlock()

	if p.closed.Get() {
		return fmt.Errorf("peer is closed")
	}

	p.routerGroup = p.routerGroupProvider.GetOrNewRouterGroup(routerGroupId)

	if p.publisherProvider != nil && !config.NoPublish {
		p.publisher = p.publisherProvider.NewPublisher(p.id, p.routerGroup)
	}

	if p.subscriberProvider != nil && !config.NoSubscribe {
		p.subscriber = p.subscriberProvider.NewSubscriber(p.id, p.routerGroup)
	}

	p.routerGroup.AddPeer(p)

	if !config.NoSubscribe {
		p.routerGroup.Subscribe(p)
	}

	return nil
}

func (p *Peer) ID() string {
	return p.id
}

func (p *Peer) RouterGroup() IRouterGroup {
	return p.routerGroup
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

	if p.routerGroup != nil {
		p.routerGroup.RemovePeer(p)
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
