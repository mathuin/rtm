package rtm

import (
	"fmt"
	"os/exec"
	"runtime"

	"golang.org/x/net/context"
)

// Session is the authenticated token associated with a client.
type Session struct {
	parent   *Client
	Frob     string
	Token    string
	Timeline *Timeline
}

// CreateSession is a required step to authenticate for further API use.
func (c *Client) CreateSession(ctx context.Context) (*Session, error) {
	c.setctx(ctx)

	s := &Session{parent: c}
	r := Request{"perms": "delete", "frob": s.frob()}
	u := c.url(c.urlAuth(), r)

	// fmt.Printf("Visit the following URL in your web browser: %s\n", u)
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", u).Start()
	case "windows", "darwin":
		err = exec.Command("open", u).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return nil, err
	}

	// time.Sleep(time.Second * 5)

	_ = s.token()
	_ = s.timeline()

	return s, nil
}

func (s *Session) frob() string {
	c := s.parent
	if s.Frob == "" {
		// Create Frob for authentication request.
		var m FrobResp

		m, err := c.Frob()
		mustNotErr(err)
		s.Frob = m.Frob
	}
	return s.Frob
}

func (s *Session) token() string {
	if s.Token == "" {
		var t TokenResp
		t, err := s.parent.Token(s.frob())
		// FIXME: check for error code 101 here and try again
		mustNotErr(err)
		s.Token = t.Auth.Token
	}
	return s.Token
}

func (s *Session) timeline() *Timeline {
	if s.Timeline == nil {
		t := Timeline{parent: s}
		_ = t.token()
		s.Timeline = &t
	}
	return s.Timeline
}

// CutTimeline sets current session timeline to "".  This renders all previous transactions un-undoable.
func (s *Session) CutTimeline() {
	s.Timeline = nil
}
