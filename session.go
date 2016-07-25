package rtm

import "golang.org/x/net/context"

// Session is the authenticated token associated with a client.
type Session struct {
	parent *Client
	Token  string
}

// Login should report what user if any is logged in
func (s *Session) Login(ctx context.Context) (LoginResp, error) {
	c := s.parent
	var m map[string]LoginResp
	r := Request{"method": "rtm.test.login", "auth_token": s.Token}
	if err := c.doReqURL(ctx, c.url(c.urlBase(), r), &m); err != nil {
		return LoginResp{}, err
	}
	lr := m["rsp"]
	if lr.Status == "fail" {
		return LoginResp{}, &lr.Error
	}
	return lr, nil
}
