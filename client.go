package gotezos

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// client is a struct to represent the http or rpc client
type client struct {
	URL       string
	netClient *http.Client
}

// newClient returns a new client
func newClient(URL string) *client {
	if URL[len(URL)-1] == '/' {
		URL = URL[:len(URL)-1]
	}
	if !strings.HasPrefix(URL, "http://") && !strings.HasPrefix(URL, "https://") {
		URL = fmt.Sprintf("http://%s", URL) //default to http
	}

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	return &client{URL: URL, netClient: netClient}
}

func (c *client) Put(path, args string) ([]byte, error) {
	var respBytes []byte
	request, err := http.NewRequest("PUT", c.URL+path, bytes.NewBuffer([]byte(args)))
	if err != nil {
		return respBytes, err
	}
	request.Header.Set("content-type", "application/json")

	return c.do(request)
}

func (c *client) Post(path, args string) ([]byte, error) {
	var respBytes []byte
	request, err := http.NewRequest("PUT", c.URL+path, bytes.NewBuffer([]byte(args)))
	if err != nil {
		return respBytes, err
	}
	request.Header.Set("content-type", "application/json")

	return c.do(request)
}

func (c *client) Get(path string, params map[string]string) ([]byte, error) {
	var bytes []byte

	req, err := http.NewRequest("GET", c.URL+path, nil)
	if err != nil {
		return bytes, err
	}

	q := req.URL.Query()
	if len(params) > 0 {
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	return c.do(req)
}

func (c *client) do(req *http.Request) ([]byte, error) {
	var bytes []byte

	resp, err := c.netClient.Do(req)
	if err != nil {
		return bytes, err
	}

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return bytes, err
	}

	if resp.StatusCode != http.StatusOK {
		return bytes, errors.Errorf("%d error: %s", resp.StatusCode, string(bytes))
	}

	c.netClient.CloseIdleConnections()

	return bytes, nil
}
