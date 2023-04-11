package forwarder

import (
	"sync"

	"github.com/let-light/neon/pkg/utils"
	"github.com/sirupsen/logrus"
)

type IGroup interface {
	ID() string
	Publish()
	AddPeer(peer IPeer)
	GetPeer(peerID string) IPeer
	RemovePeer(peer IPeer)
	Subscribe(p IPeer)
}

type Group struct {
	gid            string
	mutex          sync.RWMutex
	peers          map[string]IPeer
	closed         utils.AtomicBool
	onCloseHandler func()
	logger         *logrus.Entry
}

func NewGroup(gid string, logger *logrus.Entry) IGroup {
	g := &Group{
		gid:    gid,
		peers:  make(map[string]IPeer),
		logger: logger,
	}

	return g
}

func (g *Group) ID() string {
	return g.gid
}

func (g *Group) AddPeer(peer IPeer) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.peers[peer.ID()] = peer
}

func (g *Group) GetPeer(peerID string) IPeer {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return g.peers[peerID]
}

func (g *Group) RemovePeer(peer IPeer) {
	peerId := peer.ID()
	g.logger.Info("RemovePeer ", "peer_id", peer.ID(), "group gid", g.gid)
	g.mutex.Lock()
	if g.peers[peerId] == peer {
		delete(g.peers, peerId)
	}
	peerCount := len(g.peers)
	g.mutex.Unlock()

	// Close group if no peers
	if peerCount == 0 {
		g.Close()
	}
}

func (g *Group) Publish() {

}

func (g *Group) Subscribe(peer IPeer) {

}

// OnClose is called when the group is closed
func (g *Group) OnClose(f func()) {
	g.onCloseHandler = f
}

func (g *Group) Close() {
	if !g.closed.Set(true) {
		return
	}

	if g.onCloseHandler != nil {
		g.onCloseHandler()
	}
}
