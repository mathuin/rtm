// +build integration

package rtm

import (
	"os"
	"strconv"
	"testing"
	"time"

	"golang.org/x/net/context"
)

var c Client
var s *Session

// Create a file named info.json and put it at the root of your project.
// That file should contain a valid API key and shared secret from RTM.
// run `go test -v --tags=integration .` to start integration tests.
func TestMain(m *testing.M) {
	c = Client{}
	err := c.readSecrets("info.json")
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	s, err = c.CreateSession(ctx)
	mustNotErr(err)
	retCode := m.Run()
	// teardown
	os.Exit(retCode)
}

func TestEcho(t *testing.T) {
	expected := "pong"
	actual, err := c.Echo(expected)
	if err != nil {
		t.Errorf("err: expected nil, got %s", err.Error())
	}
	if actual.Ping != expected {
		t.Errorf("ping: expected %s, got %s", expected, actual.Ping)
	}
}

func TestLogin(t *testing.T) {
	actual, err := s.Login()
	if err != nil {
		t.Errorf("err: expected nil, got %s", err.Error())
	}
	if actual.User.ID == "" {
		t.Errorf("User.ID: expected not \"\", got \"\"")
	}
}

func TestCheckToken(t *testing.T) {
	actual, err := s.CheckToken()
	if err != nil {
		t.Errorf("err: expected nil, got %s", err.Error())
	}
	if actual.Auth.Token != s.Token {
		t.Errorf("Token: expected %s, got %s", s.Token, actual.Auth.Token)
	}
}

func TestCreateTimeline(t *testing.T) {
	actual, err := s.CreateTimeline()
	if err != nil {
		t.Errorf("err: expected nil, got %s", err.Error())
	}
	if _, err := strconv.Atoi(actual.Timeline); err != nil {
		t.Errorf("Timeline: expected integer, got %s", actual.Timeline)
	}
}
