package anh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WellKnownTag_String(t *testing.T) {
	tag := WellKnownTag{"InstallationId", "123"}
	assert.Equal(t, "$InstallationId:{123}", tag.String())
}

func Test_ParseWellKnownTag(t *testing.T) {
	tag, found := ParseWellKnownTag("$InstallationId:{123}")
	assert.True(t, found)
	assert.Equal(t, WellKnownTag{"InstallationId", "123"}, tag)

	_, found = ParseWellKnownTag("$InstallationId:123}")
	assert.False(t, found)

	_, found = ParseWellKnownTag("$InstallationId:123")
	assert.False(t, found)

	_, found = ParseWellKnownTag("InstallationId:123")
	assert.False(t, found)

	_, found = ParseWellKnownTag("InstallationId:{123}")
	assert.False(t, found)
}
