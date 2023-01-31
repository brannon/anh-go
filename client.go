package anh

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Client struct {
	HubName string
	Logger  *log.Logger

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

	if c.Logger != nil {
		c.Logger.Printf(">> URL: %s\n", req.URL.String())
		c.Logger.Printf(">> Headers:\n")
		for name := range req.Header {
			c.Logger.Printf(">>   %s: %s\n", name, req.Header.Get(name))
		}
	}

	res, err := http.DefaultClient.Do(req)

	if c.Logger != nil {
		if err != nil {
			c.Logger.Printf("<< Client error: %s\n", err.Error())
		} else {
			c.Logger.Printf("<< Status: %s\n", res.Status)
			c.Logger.Printf("<< Headers:\n")
			for name := range res.Header {
				c.Logger.Printf("<<   %s: %s\n", name, res.Header.Get(name))
			}
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

	switch res.StatusCode {
	case http.StatusOK:
		return nil

	case http.StatusUnauthorized:
		return ErrInvalidCredentials

	default:
		return fmt.Errorf("unexpected response (%s)", res.Status)
	}
}
