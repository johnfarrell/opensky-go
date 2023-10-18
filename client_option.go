package opensky

import (
	"fmt"
	"time"
)

type ClientOption interface {
	apply(*Client) error
}

type Authentication struct {
	ClientOption
	username string
	password string
}

func WithAuthentication(username, password string) *Authentication {
	return &Authentication{
		username: username,
		password: password,
	}
}

func (opt *Authentication) apply(client *Client) error {
	client.authHeader = fmt.Sprintf("%s:%s", opt.username, opt.password)
	return nil
}

type Timeout struct {
	ClientOption
	timeout time.Duration
}

func WithTimeout(timeout time.Duration) *Timeout {
	return &Timeout{
		timeout: timeout,
	}
}

func (opt *Timeout) apply(client *Client) error {
	client.httpClient.Timeout = opt.timeout
	return nil
}
