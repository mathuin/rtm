package rtm

import "golang.org/x/net/context"

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

// Frob should return a frob value.
func (c *Client) Frob(ctx context.Context) (FrobResp, error) {
	var m frobResp
	r := Request{"method": "rtm.auth.getFrob"}
	if err := c.doReqURL(ctx, c.url(c.urlBase(), r), &m); err != nil {
		return FrobResp{}, err
	}
	return m.RSP, m.RSP.IsOK()
}

// Token returns an auth Token
func (c *Client) Token(ctx context.Context, f string) (TokenResp, error) {
	var m tokenResp
	r := Request{"method": "rtm.auth.getToken", "frob": f}
	if err := c.doReqURL(ctx, c.url(c.urlBase(), r), &m); err != nil {
		return TokenResp{}, err
	}
	return m.RSP, m.RSP.IsOK()
}
