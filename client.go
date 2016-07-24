package rtm

// Based heavily on cep21/smitego

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"sort"

	"golang.org/x/net/context"
)

// AuthServicesURL is the authentication services URL.
const AuthServicesURL = "https://www.rememberthemilk.com/services/auth/"

// RESTEndpointURL is the REST endpoint URL.
const RESTEndpointURL = "https://api.rememberthemilk.com/services/rest/"

// Secrets indicate information which will be loaded from a file
type Secrets struct {
	APIKey string `json:"api_key"`
	Secret string `json:"secret"`
}

// Client can create RTM sessions and interact with the API.
type Client struct {
	APIKey  string
	Secret  string
	AuthURL string
	BaseURL string
}

func (c *Client) readSecrets(filename string) error {
	// unmarshal json from file
	sfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var s Secrets
	if err = json.Unmarshal(sfile, &s); err != nil {
		return err
	}
	c.APIKey = s.APIKey
	c.Secret = s.Secret
	return nil
}

func (c *Client) urlAuth() string {
	if c.AuthURL == "" {
		return AuthServicesURL
	}
	return c.AuthURL
}

func (c *Client) apiKey() string {
	if c.APIKey == "" {
		panic("API key undefined")
	}
	return c.APIKey
}

func (c *Client) secret() string {
	if c.Secret == "" {
		panic("Secret undefined")
	}
	return c.Secret
}

// Request contains parameters
type Request struct {
	Method     string
	Parameters map[string]string
}

// Sign will generate API signature for a particular Request
func (c *Client) Sign(r Request) string {
	// add client's APIKey as value with key 'api_key' to request parameters
	params := r.Parameters
	if r.Method != "" {
		params["method"] = r.Method
	}
	params["api_key"] = c.apiKey()

	// sort parameters by key value alphabetically
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// make a big string (but prepend shared secret)
	var buffer bytes.Buffer
	buffer.WriteString(c.Secret)
	for _, k := range keys {
		buffer.WriteString(k)
		buffer.WriteString(params[k])
	}

	// compute MD5 Sum
	hasher := md5.New()
	hasher.Write(buffer.Bytes())
	return string(hex.EncodeToString(hasher.Sum(nil)))
}

// CreateSession is a required step to authenticate for further API use.
func (c *Client) CreateSession(ctx context.Context) (*Session, error) {
	var SessionID string

	// add stuff
	r := Request{}

	r.Parameters["perms"] = "delete"
	// r.Parameters["api_sig"] = c.Sign(r)
	// u := c.url(c.authUrl(), r)

	return &Session{
		parent:    c,
		SessionID: SessionID,
	}, nil
}

func (c *Client) urlBase() string {
	if c.BaseURL == "" {
		return RESTEndpointURL
	}
	return c.BaseURL
}

func (c *Client) url(s string, r Request) string {
	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()

	// starts with method
	if r.Method != "" {
		r.Parameters["method"] = r.Method
	}

	// then api_key
	q.Set("api_key", c.apiKey())

	// then parameters
	for k, v := range r.Parameters {
		q.Set(k, v)
	}
	q.Set("api_sig", c.Sign(r))

	// then api_sig
	ret := ""
	return ret
}

func mustNotErr(err error) {
	if err != nil {
		panic("Unexpected error: " + err.Error())
	}
}
