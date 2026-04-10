package util

import (
	"bytes"
	"encoding/json"
	"io"
)

func Serialize(payload map[string]any) (*bytes.Buffer, error) {
	serializedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(serializedPayload), nil
}

func Deserialize(body io.ReadCloser) (map[string]any, error) {
	responseBuffer, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	responseBody := make(map[string]any)
	json.Unmarshal(responseBuffer, &responseBody)

	return responseBody, nil
}
