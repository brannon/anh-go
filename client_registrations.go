package anh

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

func (c *Client) GetRegistration(ctx context.Context, registrationId string) (Registration, error) {
	url := c.buildUrl(c.HubName, "registrations", registrationId)

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

	entry, err := parseAtomEntry(res.Body)
	if err != nil {
		return nil, err
	}

	if entry.Content != nil {
		if entry.Content.AppleRegistrationDescription != nil {
			return entry.Content.AppleRegistrationDescription, nil
		} else if entry.Content.GcmRegistrationDescription != nil {
			return entry.Content.GcmRegistrationDescription, nil
		}
	}

	return nil, errors.New("invalid response from server: atom entry content is missing")
}

func (c *Client) ListRegistrations(ctx context.Context) (*Collection[Registration], error) {
	items, continuationToken, err := c.fetchNextRegistrationsPage(ctx, "", 1)
	if err != nil {
		return nil, err
	}

	return &Collection[Registration]{
		items:         items,
		fetchToken:    continuationToken,
		fetchNextPage: c.fetchNextRegistrationsPage,
		pageSize:      1,
	}, nil
}

func (c *Client) fetchNextRegistrationsPage(ctx context.Context, continuationToken string, count int) ([]Registration, string, error) {
	url := c.buildUrl(c.HubName, "registrations")

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	if continuationToken != "" {
		appendQueryString(req.URL, "ContinuationToken", continuationToken)
	}

	if count > 0 {
		appendQueryString(req.URL, "$Top", strconv.FormatInt(int64(count), 10))
	}

	res, err := c.executeRequest(req)
	if err != nil {
		return nil, "", err
	}

	if res.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected error (%s)", res.Status)
	}

	continuationToken = res.Header.Get("X-MS-ContinuationToken")

	defer res.Body.Close()

	feed, err := parseAtomFeed(res.Body)
	if err != nil {
		return nil, "", err
	}

	registrations := []Registration{}

	for _, entry := range feed.Entries {
		if entry.Content != nil {
			if entry.Content.AppleRegistrationDescription != nil {
				registrations = append(registrations, entry.Content.AppleRegistrationDescription)
			} else if entry.Content.GcmRegistrationDescription != nil {
				registrations = append(registrations, entry.Content.GcmRegistrationDescription)
			}
		}
	}

	return registrations, continuationToken, nil
}
