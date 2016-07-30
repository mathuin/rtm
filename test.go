package rtm

import "golang.org/x/net/context"

// Echo should echo the sent values
func (c *Client) Echo(ctx context.Context, p string) (EchoResp, error) {
	var m echoResp
	r := Request{"method": "rtm.test.echo", "ping": p}
	if err := c.doReqURL(ctx, c.url(c.urlBase(), r), &m); err != nil {
		return EchoResp{}, err
	}
	return m.RSP, m.RSP.IsOK()
}

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
