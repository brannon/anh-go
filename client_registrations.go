package anh

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

func (c *Client) GetRegistration(ctx context.Context, registrationId string) (Registration, error) {
	url := c.buildUrl(c.hubName, "registrations", registrationId)

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

	registration := entry.GetRegistration()
	if registration == nil {
		return nil, errors.New("invalid response from server: atom entry did not contain a recognized registration")
	}

	return registration, nil
}

func (c *Client) ListRegistrations(ctx context.Context) (*PagedCollection[Registration], error) {
	items, continuationToken, err := c.fetchNextRegistrationsPage(ctx, "", 1)
	if err != nil {
		return nil, err
	}

	return &PagedCollection[Registration]{
		items:         items,
		fetchToken:    continuationToken,
		fetchNextPage: c.fetchNextRegistrationsPage,
		pageSize:      1,
	}, nil
}

func (c *Client) fetchNextRegistrationsPage(ctx context.Context, continuationToken string, count int) ([]Registration, string, error) {
	url := c.buildUrl(c.hubName, "registrations")

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
		registration := entry.GetRegistration()
		if registration != nil {
			registrations = append(registrations, registration)
		}
	}

	return registrations, continuationToken, nil
}
