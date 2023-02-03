package anh

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetInstallation(ctx context.Context, id string) (*Installation, error) {
	url := c.buildUrl(c.hubName, "installations", id)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	res, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("unexpected error (%s)", res.Status)
	}

	defer res.Body.Close()

	rawData, err := parseJSONObject(res.Body)
	if err != nil {
		return nil, err
	}

	return &Installation{rawData: rawData}, nil
}
