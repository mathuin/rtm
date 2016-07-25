// +build integration

package rtm

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

var c Client
var s *Session
var ctx context.Context
var cancel context.CancelFunc

// Create a file named info.json and put it at the root of your project.  That file should contain a valid API key and shared secret from RTM.
// run `go test -v --tags=integration .` to start integration tests.
func TestMain(m *testing.M) {
	c = Client{}
	err := c.readSecrets("info.json")
	if err != nil {
		panic(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	s, err = c.CreateSession(ctx)
	mustNotErr(err)
	retCode := m.Run()
	// teardown
	os.Exit(retCode)
}

func TestEcho(t *testing.T) {
	Convey("Echo should work", t, func() {
		expected := "pong"
		actual, err := c.Echo(ctx, expected)
		So(err, ShouldBeNil)
		So(actual.Ping, ShouldEqual, expected)
	})
}

func TestLogin(t *testing.T) {
	Convey("Login should work", t, func() {
		actual, err := s.Login(ctx)
		So(err, ShouldBeNil)
		So(actual.User.ID, ShouldNotBeNil)
	})
}
