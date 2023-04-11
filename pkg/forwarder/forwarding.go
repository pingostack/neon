package forwarder

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type forwardingSys struct {
	groups map[string]IGroup
	sync.RWMutex
}

var ForwardingSys = &forwardingSys{
	groups: make(map[string]IGroup),
}

func (f *forwardingSys) newGroup(id string) IGroup {
	group := NewGroup(id, logrus.WithField("group", id)).(*Group)

	group.OnClose(func() {
		f.removeGroup(id)
	})

	f.Lock()
	f.groups[id] = group
	f.Unlock()

	return group
}

func (f *forwardingSys) getGroup(id string) IGroup {
	f.RLock()
	defer f.RUnlock()
	return f.groups[id]
}

func (f *forwardingSys) GetOrNewGroup(sid string) IGroup {
	group := f.getGroup(sid)
	if group == nil {
		group = f.newGroup(sid)
	}
	return group
}

func (f *forwardingSys) removeGroup(gid string) {
	f.Lock()
	defer f.Unlock()

	delete(f.groups, gid)
}
