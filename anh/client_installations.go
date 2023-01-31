package anh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetInstallation(ctx context.Context, id string) (*Installation, error) {
	url := c.buildUrl(c.HubName, "installations", id)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	res, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusOK:
		defer res.Body.Close()
		rawData := make(JSONObject)
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&rawData)
		if err != nil {
			return nil, err
		}

		return &Installation{rawData: rawData}, nil

	default:
		return nil, fmt.Errorf("unexpected error (%s)", res.Status)
	}
}
