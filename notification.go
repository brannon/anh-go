package anh

import (
	"encoding/json"
	"io"
	"time"
)

type Notification interface {
	GetBody() ([]byte, error)
	GetContentType() string
	GetHeaders() map[string]string
	GetPlatform() Platform
}

type AppleNotification struct {
	Body     JSONObject
	Expiry   time.Time
	Priority int
}

func (n *AppleNotification) GetBody() ([]byte, error) {
	if n.Body == nil || len(n.Body) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(n.Body)
}

func (n *AppleNotification) GetContentType() string {
	return "application/json"
}

func (n *AppleNotification) GetHeaders() map[string]string {
	return map[string]string{}
}

func (n *AppleNotification) GetPlatform() Platform {
	return PlatformApple
}

func (n *AppleNotification) SetBodyFromReader(r io.Reader) error {
	return json.NewDecoder(r).Decode(&n.Body)
}

func (n *AppleNotification) SetBodyFromString(s string) error {
	return json.Unmarshal([]byte(s), &n.Body)
}
