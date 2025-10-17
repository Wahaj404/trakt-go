package util

import (
	"encoding/json"
	"io"
	"net/http"
)

func DeserializeResponse(resp *http.Response) (map[string]any, error) {
	defer resp.Body.Close()
	responseBuffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	responseBody := make(map[string]any)
	json.Unmarshal(responseBuffer, &responseBody)

	return responseBody, nil
}
