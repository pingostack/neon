package forwarding

import (
	"sync"

	"github.com/let-light/neon/pkg/utils"
	"github.com/sirupsen/logrus"
)

type IRouterGroup interface {
	ID() string
	Publish(router Router)
	AddPeer(peer IPeer)
	GetPeer(peerID string) IPeer
	RemovePeer(peer IPeer)
	Subscribe(p IPeer)
}

type RouterGroup struct {
	id             string
	mu             sync.RWMutex
	peers          map[string]IPeer
	closed         utils.AtomicBool
	onCloseHandler func()
	logger         *logrus.Entry
}

func NewRouterGroup() *RouterGroup {
	return &RouterGroup{}
}

func (g *RouterGroup) ID() string {
	return g.id
}

func (g *RouterGroup) AddPeer(peer IPeer) {
	g.mu.Unlock()
	defer g.mu.Unlock()

	g.peers[peer.ID()] = peer
}

func (g *RouterGroup) GetPeer(peerID string) IPeer {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.peers[peerID]
}

func (g *RouterGroup) RemovePeer(peer IPeer) {
	peerId := peer.ID()
	g.logger.Info("RemovePeer ", "peer_id", peer.ID(), "routerGroup id", g.id)
	g.mu.Lock()
	if g.peers[peerId] == peer {
		delete(g.peers, peerId)
	}
	peerCount := len(g.peers)
	g.mu.Unlock()

	// Close routerGroup if no peers
	if peerCount == 0 {
		g.Close()
	}
}

func (g *RouterGroup) Publish(peer IPeer) {

}

func (g *RouterGroup) Subscribe(peer IPeer) {

}

// OnClose is called when the routerGroup is closed
func (g *RouterGroup) OnClose(f func()) {
	g.onCloseHandler = f
}

func (g *RouterGroup) Close() {
	if !g.closed.Set(true) {
		return
	}

	if g.onCloseHandler != nil {
		g.onCloseHandler()
	}
}
