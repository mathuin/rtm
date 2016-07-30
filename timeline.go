package rtm

// Timeline is a set of sub-transactions which can be undone.
type Timeline struct {
	parent       *Session
	Token        string
	Transactions []string
}

func (t *Timeline) token() string {
	s := t.parent
	if t.Token == "" {
		timeline, err := s.CreateTimeline()
		mustNotErr(err)
		t.Token = timeline.Timeline
	}
	return t.Token
}

// CreateTimeline will return a newly created timeline.
func (s *Session) CreateTimeline() (TimelineResp, error) {
	c := s.parent
	var m timelineResp
	r := Request{"method": "rtm.timelines.create", "auth_token": s.Token}
	if err := c.doReqURL(c.url(c.urlBase(), r), &m); err != nil {
		return TimelineResp{}, err
	}
	return m.RSP, m.RSP.IsOK()
}
