package anh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONObject_String(t *testing.T) {
	var obj JSONObject

	obj = JSONObject{
		"foo": "bar",
	}
	assert.Equal(t, `{"foo":"bar"}`, obj.String())

	obj = JSONObject{
		"foo": "bar",
		"baz": 123,
	}
	assert.Equal(t, `{"baz":123,"foo":"bar"}`, obj.String())
}

func TestJSONObject_FormattedString(t *testing.T) {
	var obj JSONObject

	obj = JSONObject{
		"foo": "bar",
	}
	assert.Equal(t, `{
  "foo": "bar"
}
`, obj.FormattedString())

	obj = JSONObject{
		"foo": "bar",
		"baz": 123,
	}
	assert.Equal(t, `{
  "baz": 123,
  "foo": "bar"
}
`, obj.FormattedString())
}
