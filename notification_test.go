package anh

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AppleNotification_GetBody(t *testing.T) {
	notification := &AppleNotification{}

	body, err := notification.GetBody()

	assert.NoError(t, err)
	assert.Equal(t, []byte("{}"), body)

	notification.Body = JSONObject{
		"aps": JSONObject{
			"alert": JSONObject{
				"title": "Hello",
				"body":  "World",
			},
		},
	}

	body, err = notification.GetBody()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"aps":{"alert":{"title":"Hello","body":"World"}}}`, string(body))
}

func Test_AppleNotification_GetContentType(t *testing.T) {
	notification := &AppleNotification{}

	assert.Equal(t, "application/json", notification.GetContentType())
}

func Test_AppleNotification_GetHeaders(t *testing.T) {
	notification := &AppleNotification{}

	headers := notification.GetHeaders()

	assert.Equal(t, 0, len(headers))
}

func Test_AppleNotification_GetPlatform(t *testing.T) {
	notification := &AppleNotification{}

	assert.Equal(t, PlatformApple, notification.GetPlatform())
}

func Test_AppleNotification_SetBodyFromReader(t *testing.T) {
	notification := &AppleNotification{}
	r := strings.NewReader(`{"aps":{"alert":{"body":"World","title":"Hello"}}}`)

	err := notification.SetBodyFromReader(r)

	assert.NoError(t, err)
	assert.JSONEq(t, `{"aps":{"alert":{"body":"World","title":"Hello"}}}`, notification.Body.String())
}

func Test_AppleNotification_SetBodyFromString(t *testing.T) {
	notification := &AppleNotification{}

	err := notification.SetBodyFromString(`{"aps":{"alert":{"body":"World","title":"Hello"}}}`)

	assert.NoError(t, err)
	assert.JSONEq(t, `{"aps":{"alert":{"body":"World","title":"Hello"}}}`, notification.Body.String())

}
