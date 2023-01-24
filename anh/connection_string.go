package anh

import (
	"strings"

	"github.com/pkg/errors"
)

type ConnectionString struct {
	Endpoint string
	KeyName  string
	Key      string
}

func ParseConnectionString(s string) (*ConnectionString, error) {
	cs := ConnectionString{}

	pairs := strings.Split(s, ";")
	for _, pair := range pairs {
		keyAndValue := strings.SplitN(pair, "=", 2)
		if len(keyAndValue) != 2 {
			return nil, errors.New("invalid connection string")
		}

		if keyAndValue[0] != "" {
			if strings.EqualFold("Endpoint", keyAndValue[0]) {
				cs.Endpoint = keyAndValue[1]
			} else if strings.EqualFold("SharedAccessKeyName", keyAndValue[0]) {
				cs.KeyName = keyAndValue[1]
			} else if strings.EqualFold("SharedAccessKey", keyAndValue[0]) {
				cs.Key = keyAndValue[1]
			}
		}
	}

	if cs.Endpoint == "" {
		return nil, errors.New("missing connection string value for 'Endpoint'")
	}
	if cs.KeyName == "" {
		return nil, errors.New("missing connection string value for 'SharedAccessKeyName'")
	}
	if cs.Key == "" {
		return nil, errors.New("missing connection string value for 'Key'")
	}

	return &cs, nil
}
