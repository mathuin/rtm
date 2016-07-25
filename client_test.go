package rtm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

// keeping this as a table-driven test example until I write another one
// Examples from http://automation.frankmontanaro.com/setting-up-remember-the-milks-api/
var signtests = []struct {
	clientKey     string
	clientSecret  string
	requestParams map[string]string
	expected      string
}{
	{"1234567890", "987654321", map[string]string{"perms": "delete"}, "efb5d96f44d33b72081b81ddde96005d"},
	{"1234567890", "987654321", map[string]string{"method": "rtm.tasks.getList", "auth_token": "666a777b999c", "format": "json", "filter": "status:incomplete AND due:never OR due:today"}, "9094df9d88641c8c0f5666accc335761"},
}

func TestSign(t *testing.T) {
	for _, tt := range signtests {
		c := Client{
			APIKey: tt.clientKey,
			Secret: tt.clientSecret,
		}
		r := tt.requestParams
		actual := c.Sign(r)
		if actual != tt.expected {
			t.Errorf("Sign with Client key %q, secret %q, and Request params %q => %q, want %q", tt.clientKey, tt.clientSecret, tt.requestParams, actual, tt.expected)
		}
	}
}

func TestReadSecrets(t *testing.T) {
	inkey := "mykey"
	insecret := "mysecret"
	s := Secrets{
		APIKey: inkey,
		Secret: insecret,
	}
	b, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	tmpfile, err := ioutil.TempFile("", "secrets")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpfile.Write(b); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	c := Client{}
	c.readSecrets(tmpfile.Name())

	if c.APIKey != inkey {
		t.Errorf("token: expected %s, got %s", inkey, c.APIKey)
	}

	if c.Secret != insecret {
		t.Errorf("token: expected %s, got %s", insecret, c.Secret)
	}
}

var urltests = []struct {
	clientKey     string
	clientSecret  string
	requestParams map[string]string
	expected      string
}{
	{"1234567890", "987654321", map[string]string{"perms": "delete"}, "https://www.rememberthemilk.com/services/auth/?api_key=1234567890&api_sig=076f380119ba2ff3a7c20f060d864533&format=json&perms=delete"},
	{"1234567890", "987654321", map[string]string{"method": "rtm.tasks.getList", "auth_token": "666a777b999c", "filter": "status:incomplete AND due:never OR due:today"}, "https://www.rememberthemilk.com/services/auth/?api_key=1234567890&api_sig=9094df9d88641c8c0f5666accc335761&auth_token=666a777b999c&filter=status%3Aincomplete+AND+due%3Anever+OR+due%3Atoday&format=json&method=rtm.tasks.getList"},
}

func TestUrl(t *testing.T) {
	for _, tt := range urltests {
		c := Client{
			APIKey: tt.clientKey,
			Secret: tt.clientSecret,
		}
		r := tt.requestParams
		s := AuthServicesURL
		actual := c.url(s, r)
		if actual != tt.expected {
			t.Errorf("url: expected %s, got %s", tt.expected, actual)
		}
	}
}
