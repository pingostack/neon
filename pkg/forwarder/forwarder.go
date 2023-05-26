package forwarder

import (
	"sync"
)

type forwarder struct {
	groups map[string]IGroup
	sync.RWMutex
}

var Forwarder = &forwarder{
	groups: make(map[string]IGroup),
}

func (f *forwarder) newGroup(id string) IGroup {
	group := NewGroup(id).(*Group)

	group.OnClose(func() {
		f.removeGroup(id)
	})

	f.Lock()
	f.groups[id] = group
	f.Unlock()

	return group
}

func (f *forwarder) getGroup(id string) IGroup {
	f.RLock()
	defer f.RUnlock()
	return f.groups[id]
}

func (f *forwarder) GetOrNewGroup(sid string) IGroup {
	group := f.getGroup(sid)
	if group == nil {
		group = f.newGroup(sid)
	}
	return group
}

func (f *forwarder) removeGroup(gid string) {
	f.Lock()
	defer f.Unlock()

	delete(f.groups, gid)
}
