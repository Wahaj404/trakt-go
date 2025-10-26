package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func SerializeRequest(payload map[string]any) (*bytes.Buffer, error) {
	serializedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(serializedPayload), nil
}

func DeserializeResponse(resp *http.Response) (map[string]any, error) {
	responseBuffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	responseBody := make(map[string]any)
	json.Unmarshal(responseBuffer, &responseBody)

	return responseBody, nil
}
