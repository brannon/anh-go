package anh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseConnectionString(t *testing.T) {
	cs, err := ParseConnectionString("Endpoint=endpoint;SharedAccessKeyName=keyname;SharedAccessKey=key")
	assert.Nil(t, err)
	assert.Equal(t, "endpoint", cs.Endpoint)
	assert.Equal(t, "keyname", cs.KeyName)
	assert.Equal(t, "key", cs.Key)
}

func Test_ParseConnectionString_Invalid(t *testing.T) {
	var err error

	_, err = ParseConnectionString("SharedAccessKeyName=keyname;SharedAccessKey=key")
	assert.ErrorContains(t, err, "Endpoint")

	_, err = ParseConnectionString("Endpoint=endpoint;SharedAccessKey=key")
	assert.ErrorContains(t, err, "SharedAccessKeyName")

	_, err = ParseConnectionString("Endpoint=endpoint;SharedAccessKeyName=keyname")
	assert.ErrorContains(t, err, "SharedAccessKey")
}
