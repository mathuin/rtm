package rtm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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

func TestSignTable(t *testing.T) {
	for _, tt := range signtests {
		c := Client{
			APIKey: tt.clientKey,
			Secret: tt.clientSecret,
		}
		r := Request{
			Parameters: tt.requestParams,
		}
		actual := c.Sign(r)
		if actual != tt.expected {
			t.Errorf("Sign with Client key %q, secret %q, and Request params %q => %q, want %q", tt.clientKey, tt.clientSecret, tt.requestParams, actual, tt.expected)
		}
	}
}

func TestSignConvey(t *testing.T) {
	for _, tt := range signtests {
		Convey("Given a client with known key and secret, and a request with known parameters", t, func() {
			c := Client{
				APIKey: tt.clientKey,
				Secret: tt.clientSecret,
			}
			r := Request{
				Parameters: tt.requestParams,
			}
			Convey("When the request is signed", func() {
				actual := c.Sign(r)
				Convey("The signature should match the expected value", func() {
					So(actual, ShouldEqual, tt.expected)
				})
			})
		})
	}
}

func TestReadSecrets(t *testing.T) {
	inkey := "mykey"
	insecret := "mysecret"
	Convey("Given a known working JSON file with key and secrets", t, func() {
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
		Convey("When the file is read", func() {
			c := Client{}
			c.readSecrets(tmpfile.Name())
			Convey("The client's secrets should be set to the values from the file", func() {
				So(c.APIKey, ShouldEqual, inkey)
				So(c.Secret, ShouldEqual, insecret)
			})
		})
	})
}
