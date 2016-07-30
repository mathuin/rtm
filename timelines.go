package rtm

import "golang.org/x/net/context"

// Timeline will return a newly created timeline.
func (s *Session) Timeline(ctx context.Context) (TimelineResp, error) {
	c := s.parent
	var m timelineResp
	r := Request{"method": "rtm.timelines.create", "auth_token": s.Token}
	if err := c.doReqURL(ctx, c.url(c.urlBase(), r), &m); err != nil {
		return TimelineResp{}, err
	}
	return m.RSP, m.RSP.IsOK()
}
