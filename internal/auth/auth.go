package auth

import (
	"net/http"
	"strings"
)

// GetAPIKey extracts the API key from the "Authorization" header in the provided HTTP headers.
// The expected format of the header is "Authorization: ApiKey <key>".
// Returns the API key as a string if present in the correct format, or an error if not.

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", http.ErrNoCookie
	}
	apiKeys := strings.Split(apiKey, " ")
	if len(apiKeys) != 2 || apiKeys[0] != "ApiKey" {
		return "", http.ErrNoCookie
	}
	return apiKeys[1], nil
}