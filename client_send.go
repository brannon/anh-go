package anh

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func parseNotificationIdFromStatusURL(u *url.URL) string {
	statusPath := u.Path
	parts := strings.Split(statusPath, "/")
	return parts[len(parts)-1]
}

type NotificationResult struct {
	statusURL *url.URL
}

func (r *NotificationResult) GetNotificationId() string {
	if r.statusURL != nil {
		return parseNotificationIdFromStatusURL(r.statusURL)
	}
	return ""
}

func (c *Client) SendDirectNotification(ctx context.Context, notification Notification, deviceToken ...string) (*NotificationResult, error) {
	requestURL := c.buildUrl(c.hubName, "messages")

	req, _ := http.NewRequestWithContext(ctx, "POST", requestURL, nil)

	appendQueryString(req.URL, "direct", "true")

	req.Header.Set("ServiceBusNotification-DeviceHandle", deviceToken[0])
	req.Header.Set("ServiceBusNotification-Format", notification.GetPlatform().String())

	req.Header.Set("Content-Type", notification.GetContentType())

	bodyData, err := notification.GetBody()
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(bodyData))

	res, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	err = c.checkResponse(res, http.StatusCreated, http.StatusOK)
	if err != nil {
		return nil, err
	}

	result := &NotificationResult{}

	statusURLString := res.Header.Get("Location")
	if statusURLString != "" {
		statusURL, err := url.Parse(statusURLString)
		if err == nil {
			result.statusURL = statusURL
		}
	}

	return result, nil
}
