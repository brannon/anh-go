package anh

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type Client struct {
	VerboseLogger io.Writer

	endpoint      string
	httpClient    *http.Client
	hubName       string
	tokenProvider TokenProvider
}

type ClientOption func(*Client) error

func WithConnectionString(connectionString string) ClientOption {
	return func(c *Client) error {
		cs, err := ParseConnectionString(connectionString)
		if err != nil {
			return err
		}

		endpoint := strings.Replace(cs.Endpoint, "sb://", "https://", 1)
		tokenProvider := NewSasTokenProvider(cs.KeyName, cs.Key)

		c.endpoint = endpoint
		c.tokenProvider = tokenProvider
		return nil
	}
}

func WithHttpClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

func WithVerboseLogger(w io.Writer) ClientOption {
	return func(c *Client) error {
		c.VerboseLogger = w
		return nil
	}
}

func NewClient(hubName string, opts ...ClientOption) (*Client, error) {
	c := &Client{
		httpClient: http.DefaultClient,
		hubName:    hubName,
	}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) buildUrl(paths ...string) string {
	result, _ := url.JoinPath(c.endpoint, paths...)
	return result
}

func (c *Client) checkResponse(res *http.Response, expectedStatuses ...int) error {
	if slices.Contains(expectedStatuses, res.StatusCode) {
		return nil
	}

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return ErrInvalidCredentials

	case http.StatusNotFound:
		return ErrNotFound

	default:
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
}

func (c *Client) executeRequest(req *http.Request) (*http.Response, error) {
	if c.tokenProvider != nil {
		sasToken, _, err := c.tokenProvider.GenerateSasToken(c.endpoint, time.Now().UTC().Add(5*time.Minute))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", sasToken)
	}

	appendQueryString(req.URL, "api-version", "2020-06")

	c.verboseLogInfo("Begin Request")
	c.verboseLogOut("%s %s", req.Method, pathForLogging(req.URL))
	// Headers don't contain 'Host' by default, so fake it for logging purposes
	c.verboseLogOut("Host: %s", req.Host)
	for name := range req.Header {
		c.verboseLogOut("%s: %s", name, req.Header.Get(name))
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		c.verboseLogError("Client error: %s", err.Error())
	} else {
		c.verboseLogInfo("Begin Response")
		c.verboseLogIn("%s", res.Status)
		for name := range res.Header {
			c.verboseLogIn("%s: %s", name, res.Header.Get(name))
		}
	}

	return res, err
}

func (c *Client) Validate(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.buildUrl(c.hubName), nil)
	if err != nil {
		return err
	}

	res, err := c.executeRequest(req)
	if err != nil {
		return errors.Wrap(err, "cannot validate connection")
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return ErrInvalidCredentials
		}
		return fmt.Errorf("unexpected response (%s)", res.Status)

	}

	return nil
}

func (c *Client) verboseLogError(format string, args ...interface{}) {
	if c.VerboseLogger != nil {
		s := fmt.Sprintf(format, args...)
		c.VerboseLogger.Write([]byte("* "))
		c.VerboseLogger.Write([]byte(s))
		c.VerboseLogger.Write([]byte("\n"))
	}
}

func (c *Client) verboseLogIn(format string, args ...interface{}) {
	if c.VerboseLogger != nil {
		s := fmt.Sprintf(format, args...)
		c.VerboseLogger.Write([]byte("< "))
		c.VerboseLogger.Write([]byte(s))
		c.VerboseLogger.Write([]byte("\n"))
	}
}

func (c *Client) verboseLogInfo(format string, args ...interface{}) {
	if c.VerboseLogger != nil {
		s := fmt.Sprintf(format, args...)
		c.VerboseLogger.Write([]byte("* "))
		c.VerboseLogger.Write([]byte(s))
		c.VerboseLogger.Write([]byte("\n"))
	}
}

func (c *Client) verboseLogOut(format string, args ...interface{}) {
	if c.VerboseLogger != nil {
		s := fmt.Sprintf(format, args...)
		c.VerboseLogger.Write([]byte("> "))
		c.VerboseLogger.Write([]byte(s))
		c.VerboseLogger.Write([]byte("\n"))
	}
}
