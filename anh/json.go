package anh

import (
	"bytes"
	"encoding/json"
)

type JSONObject map[string]interface{}

func (obj JSONObject) String() string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}

func (obj JSONObject) PrettyString() string {
	bytes := bytes.Buffer{}
	encoder := json.NewEncoder(&bytes)
	encoder.SetIndent("", "  ")

	encoder.Encode(obj)

	return bytes.String()
}
