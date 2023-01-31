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
)

type Client struct {
	HubName       string
	VerboseLogger io.Writer

	endpoint      string
	tokenProvider TokenProvider
}

func NewClient(hubName string, connectionString string) (*Client, error) {
	cs, err := ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	endpoint := strings.Replace(cs.Endpoint, "sb://", "https://", 1)
	tokenProvider := NewSasTokenProvider(cs.KeyName, cs.Key)

	return &Client{
		HubName: hubName,

		endpoint:      endpoint,
		tokenProvider: tokenProvider,
	}, nil
}

func (c *Client) buildUrl(paths ...string) string {
	result, _ := url.JoinPath(c.endpoint, paths...)
	return result
}

func (c *Client) executeRequest(req *http.Request) (*http.Response, error) {
	sasToken, _, err := c.tokenProvider.GenerateSasToken(c.endpoint, time.Now().UTC().Add(5*time.Minute))
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("api-version", "2020-06")
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Authorization", sasToken)

	c.verboseLogInfo("Begin Request")
	c.verboseLogOut("%s %s", req.Method, pathForLogging(req.URL))
	// Headers don't contain 'Host' by default, so fake it for logging purposes
	c.verboseLogOut("Host: %s", req.Host)
	for name := range req.Header {
		c.verboseLogOut("%s: %s", name, req.Header.Get(name))
	}

	res, err := http.DefaultClient.Do(req)

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
	req, err := http.NewRequestWithContext(ctx, "GET", c.buildUrl(c.HubName), nil)
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
