package anh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RegistrationDescription_AddTag(t *testing.T) {
	desc := &RegistrationDescription{}
	desc.AddTag("tag1")
	assert.Equal(t, "tag1", desc.Tags)

	desc.AddTag("tag2")
	assert.Equal(t, "tag1,tag2", desc.Tags)
}

func Test_RegistrationDescription_GetTags(t *testing.T) {
	desc := &RegistrationDescription{}
	desc.Tags = "tag1,tag2"
	assert.Equal(t, []string{"tag1", "tag2"}, desc.GetTags())
}

func Test_RegistrationDescription_RemoveTag(t *testing.T) {
	desc := &RegistrationDescription{}
	desc.Tags = "tag1,tag2"
	desc.RemoveTag("tag1")
	assert.Equal(t, "tag2", desc.Tags)

	desc.RemoveTag("tag2")
	assert.Equal(t, "", desc.Tags)
}

func Test_RegistrationDescription_SetTags(t *testing.T) {
	desc := &RegistrationDescription{}
	desc.SetTags([]string{"tag1", "tag2"})
	assert.Equal(t, "tag1,tag2", desc.Tags)

	desc.SetTags([]string{})
	assert.Equal(t, "", desc.Tags)
}
