package rtm

import "fmt"

// ErrorResp is the expected response when errors occur
type ErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

func (e *ErrorResp) Error() string {
	return fmt.Sprintf("code %s, message %s\n", e.Code, e.Message)
}

type baseResp struct {
	Status string    `json:"stat"`
	Error  ErrorResp `json:"err"`
}

// FrobResp is the expected response from rtm.auth.getFrob
type FrobResp struct {
	baseResp
	Frob string `json:"frob"`
}

// TokenResp is the expected response from rtm.auth.getToken
type TokenResp struct {
	baseResp
	Auth AuthResp `json:"auth"`
}

// AuthResp is the content of the auth tag from rtm.auth.getToken
type AuthResp struct {
	Token string `json:"token"`
	Perms string `json:"perms"`
	// User
}

// EchoResp is the expected response from rtm.test.echo
type EchoResp struct {
	baseResp
	Ping string `json:"ping"`
}

// LoginResp is the expected response from rtm.test.login
type LoginResp struct {
	baseResp
}

type arbResp interface{}
