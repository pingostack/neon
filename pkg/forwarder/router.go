package forwarder

type IFilter interface {
}

type IRouter interface {
	ID() string
	AddDownTrack(track IDownTrack) error
	AddUpTrack(track IUpTrack) error
}

type Router struct {
	id    string
	group IGroup
}

func NewRouter(id string, group IGroup) IRouter {
	return &Router{
		id:    id,
		group: group,
	}
}

func (r *Router) ID() string {
	return r.id
}

func (r *Router) AddDownTrack(track IDownTrack) error {
	return nil
}

func (r *Router) AddUpTrack(track IUpTrack) error {
	return nil
}
