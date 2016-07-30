package rtm

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/net/context"
)

// Session is the authenticated token associated with a client.
type Session struct {
	parent *Client
	Token  string
}

// CreateSession is a required step to authenticate for further API use.
func (c *Client) CreateSession(ctx context.Context) (*Session, error) {
	var m FrobResp

	m, err := c.Frob(ctx)
	if err != nil {
		return nil, err
	}
	r := Request{"perms": "delete", "frob": m.Frob}
	u := c.url(c.urlAuth(), r)

	// fmt.Printf("Visit the following URL in your web browser: %s\n", u)
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

	time.Sleep(time.Second * 5)
	var t TokenResp
	t, err = c.Token(ctx, m.Frob)
	// FIXME: check for error code 101 here and try again
	if err != nil {
		return nil, err
	}

	return &Session{
		parent: c,
		Token:  t.Auth.Token,
	}, nil
}
