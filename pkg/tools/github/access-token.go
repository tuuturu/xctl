package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func requestAccessToken(client http.Client, clientID string, deviceCode string) (string, error) {
	payload, err := json.Marshal(accessTokenRequest{ClientID: clientID, DeviceCode: deviceCode, GrantType: grantType})
	if err != nil {
		return "", fmt.Errorf("preparing payload: %w", err)
	}

	request, err := http.NewRequest(
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		bytes.NewReader(payload),
	)
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	rawResponse, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}

	defer func() {
		_ = rawResponse.Body.Close()
	}()

	body, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		return "", fmt.Errorf("buffering: %w", err)
	}

	err = asError(body)
	if err != nil {
		return "", err
	}

	var response accessTokenResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("unmarshalling response: %w", err)
	}

	return response.AccessToken, nil
}
