package rtm

import "golang.org/x/net/context"

// Session is the authenticated token associated with a client.
type Session struct {
	parent *Client
	Token  string
}

// All authenticated commands should be based off session.
// Unauthenticated commands can be based off client.

// Login reports what user if any is logged in.
func (s *Session) Login(ctx context.Context) (LoginResp, error) {
	c := s.parent
	var m loginResp
	r := Request{"method": "rtm.test.login", "auth_token": s.Token}
	if err := c.doReqURL(ctx, c.url(c.urlBase(), r), &m); err != nil {
		return LoginResp{}, err
	}
	return m.RSP, m.RSP.IsOK()
}

// CheckToken will check validity of supplied token.
func (s *Session) CheckToken(ctx context.Context) (TokenResp, error) {
	c := s.parent
	var m tokenResp
	r := Request{"method": "rtm.auth.checkToken", "auth_token": s.Token}
	if err := c.doReqURL(ctx, c.url(c.urlBase(), r), &m); err != nil {
		return TokenResp{}, err
	}
	return m.RSP, m.RSP.IsOK()
}
