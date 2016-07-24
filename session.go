package rtm

// Session is the authenticated token associated with a client.
type Session struct {
	parent *Client
	Token  string
}
