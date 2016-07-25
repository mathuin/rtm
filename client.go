package rtm

// Based heavily on cep21/smitego

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"sort"
	"time"

	"golang.org/x/net/context"
)

// AuthServicesURL is the authentication services URL.
const AuthServicesURL = "https://www.rememberthemilk.com/services/auth/"

// RESTEndpointURL is the REST endpoint URL.
const RESTEndpointURL = "https://api.rememberthemilk.com/services/rest/"

// Version is the version of the library.
const Version = "0.1.0"

// Secrets indicate information which will be loaded from a file
type Secrets struct {
	APIKey string `json:"api_key"`
	Secret string `json:"secret"`
}

// Client can create RTM sessions and interact with the API.
type Client struct {
	APIKey     string
	Secret     string
	HTTPClient http.Client
	AuthURL    string
	BaseURL    string
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

func (c *Client) urlAuth() string {
	if c.AuthURL == "" {
		return AuthServicesURL
	}
	return c.AuthURL
}

func (c *Client) urlBase() string {
	if c.BaseURL == "" {
		return RESTEndpointURL
	}
	return c.BaseURL
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

// Request contains parameters
type Request map[string]string

// Sign will generate API signature for a particular Request
func (c *Client) Sign(r Request) string {
	// add client's APIKey as value with key 'api_key' to request parameters
	params := r
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

func (c *Client) url(s string, r Request) string {
	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()

	// force format to json
	r["format"] = "json"

	// then api_key
	q.Set("api_key", c.apiKey())

	// then parameters
	for k, v := range r {
		q.Set(k, v)
	}

	// then api_sig
	q.Set("api_sig", c.Sign(r))

	// Put it all together
	u.RawQuery = q.Encode()

	// Now return the string!
	return u.String()
}

func (c *Client) doReqURL(ctx context.Context, u string, jsonInto interface{}) error {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("github.com/mathuin/rtm/%s (gover %s)", Version, runtime.Version()))
	req.Cancel = ctx.Done()
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	if _, err := io.Copy(&b, resp.Body); err != nil {
		return err
	}
	debug := b.String()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %q %q", resp.StatusCode, debug)
	}
	if err := json.NewDecoder(&b).Decode(jsonInto); err != nil {
		return fmt.Errorf("not expected: %q %q", err, debug)
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	return nil
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

// Echo should echo the sent values
func (c *Client) Echo(ctx context.Context, p string) (EchoResp, error) {
	var m echoResp
	r := Request{"method": "rtm.test.echo", "ping": p}
	if err := c.doReqURL(ctx, c.url(c.urlBase(), r), &m); err != nil {
		return EchoResp{}, err
	}
	return m.RSP, m.RSP.IsOK()
}

func mustNotErr(err error) {
	if err != nil {
		panic("Unexpected error: " + err.Error())
	}
}
