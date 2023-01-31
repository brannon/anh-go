package anh

import "time"

type Installation struct {
	// Platform           string    `json:"platform"`
	// PushChannel        string    `json:"pushChannel"`
	// InstallationId     string    `json:"installationId"`
	// PushChannelExpired bool      `json:"pushChannelExpired"`
	// ExpirationTime     time.Time `json:"expirationTime"`
	// Tags               []string  `json:"tags"`

	rawData JSONObject
}

func (i *Installation) getBoolValue(name string) bool {
	if value, found := i.rawData[name]; found {
		if boolValue, ok := value.(bool); ok {
			return boolValue
		}
	}
	return false
}

func (i *Installation) getStringValue(name string) string {
	if value, found := i.rawData[name]; found {
		if stringValue, ok := value.(string); ok {
			return stringValue
		}
	}
	return ""
}

func (i *Installation) getTimeValue(name string) time.Time {
	stringValue := i.getStringValue(name)
	if stringValue != "" {
		timeValue, err := time.Parse(time.RFC3339, stringValue)
		if err == nil {
			return timeValue
		}
	}
	return time.Time{}
}

func (i *Installation) ExpirationTime() time.Time {
	return i.getTimeValue("expirationTime")
}

func (i *Installation) GetRawData() JSONObject {
	return i.rawData
}

func (i *Installation) InstallationId() string {
	return i.getStringValue("installationId")
}

func (i *Installation) Platform() string {
	return i.getStringValue("platform")
}

func (i *Installation) PushChannel() string {
	return i.getStringValue("pushChannel")
}

func (i *Installation) PushChannelExpired() bool {
	return i.getBoolValue("pushChannelExpired")
}

func (i *Installation) Tags() []string {
	if value, found := i.rawData["tags"]; found {
		if arrayValue, ok := value.([]interface{}); ok {
			stringArray := []string{}
			for _, v := range arrayValue {
				if stringValue, ok := v.(string); ok {
					stringArray = append(stringArray, stringValue)
				}
			}
			return stringArray
		}
	}
	return []string{}
}
