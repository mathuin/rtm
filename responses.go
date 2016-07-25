package rtm

import "fmt"

type errorResp struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

func (e *errorResp) Error() string {
	return fmt.Sprintf("code %s, message %s\n", e.Code, e.Message)
}

type baseResp struct {
	Status string    `json:"stat"`
	Error  errorResp `json:"err"`
}

func (b *baseResp) IsOK() error {
	if b.Status == "fail" {
		return &b.Error
	}
	return nil
}

// FrobResp is the expected response from rtm.auth.getFrob
type FrobResp struct {
	baseResp
	Frob string `json:"frob"`
}

type frobResp struct {
	RSP FrobResp `json:"rsp"`
}

// TokenResp is the expected response from rtm.auth.getToken
type TokenResp struct {
	baseResp
	Auth struct {
		Token string `json:"token"`
		Perms string `json:"perms"`
		User  struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Fullname string `json:"fullname"`
		} `json:"user"`
	} `json:"auth"`
}

type tokenResp struct {
	RSP TokenResp `json:"rsp"`
}

// EchoResp is the expected response from rtm.test.echo
type EchoResp struct {
	baseResp
	Ping string `json:"ping"`
}

type echoResp struct {
	RSP EchoResp `json:"rsp"`
}

// LoginResp is the expected response from rtm.test.login
type LoginResp struct {
	baseResp
	User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
}

type loginResp struct {
	RSP LoginResp `json:"rsp"`
}

type arbResp interface{}
