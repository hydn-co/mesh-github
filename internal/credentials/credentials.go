package credentials

import (
	"encoding/json"
	"fmt"

	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
)

type tokenCredentials struct {
	Token string `json:"token"`
}

// ExtractToken returns the GitHub personal access token from the standard api_key envelope,
// while keeping a legacy token fallback for existing stored credentials.
func ExtractToken(raw json.RawMessage) (string, error) {
	token, err := connectorutil.ExtractAPIKey(raw)
	if err == nil {
		return token, nil
	}

	if len(raw) == 0 {
		return "", err
	}

	var legacy tokenCredentials
	if decodeErr := json.Unmarshal(raw, &legacy); decodeErr != nil {
		return "", err
	}

	legacyToken, normalizeErr := connectorutil.NormalizeToken(legacy.Token)
	if normalizeErr == nil {
		return legacyToken, nil
	}

	return "", fmt.Errorf("extract token: %w", err)
}
