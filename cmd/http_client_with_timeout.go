package cmd

import (
	"errors"
	"net/http"
	"time"
)

type FuncClientOpt func(*ClientOpt)

type ClientOpt struct {
	timeOut  time.Duration
	checkRed func(req *http.Request, via []*http.Request) error
}

// defaultClientOpt initializes and returns a default ClientOpt
// instance with a timeout of 60 seconds and no redirect policy.
func defaultClientOpt() *ClientOpt {
	return &ClientOpt{
		timeOut:  60 * time.Second,
		checkRed: nil,
	}
}

func withTimeOut(d time.Duration) FuncClientOpt {
	return func(co *ClientOpt) {
		co.timeOut = d
	}
}

func witchRedirect(co *ClientOpt) {
	co.checkRed = redirectPolicyFunc
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	if len(via) >= 1 {
		return errors.New("stopped after 1 redirect")
	}
	return nil
}

// createHTTPClientWithTimeout returns a new HTTP client with the specified options.
// The function applies the options in order and returns the created client.
//
// The available options are:
// - withTimeOut(d time.Duration): sets the timeout for requests.
// - witchRedirect(): sets the redirect policy.
func createHTTPClientWithTimeout(co ...FuncClientOpt) *http.Client {
	o := defaultClientOpt()

	for _, fn := range co {
		fn(o)
	}
	client := http.Client{
		Timeout:       o.timeOut,
		CheckRedirect: o.checkRed,
	}
	return &client
}
