package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RequestDeviceCode knows how to retrieve the device code for the Device Code Flow
func RequestDeviceCode(client http.Client, clientID string) (DeviceCodeResponse, error) {
	payload, err := json.Marshal(deviceCodeRequest{ClientID: clientID})
	if err != nil {
		return DeviceCodeResponse{}, fmt.Errorf("preparing payload: %w", err)
	}

	request, err := http.NewRequest(
		http.MethodPost,
		"https://github.com/login/device/code",
		bytes.NewReader(payload),
	)
	if err != nil {
		return DeviceCodeResponse{}, fmt.Errorf("building request: %w", err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	rawResponse, err := client.Do(request)
	if err != nil {
		return DeviceCodeResponse{}, fmt.Errorf("executing request: %w", err)
	}

	rawResponseBytes, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		return DeviceCodeResponse{}, fmt.Errorf("buffering: %w", err)
	}

	var response DeviceCodeResponse

	err = json.Unmarshal(rawResponseBytes, &response)
	if err != nil {
		return DeviceCodeResponse{}, fmt.Errorf("unmarshalling: %w", err)
	}

	return response, nil
}

// PollForAccessToken knows how to poll and wait for a response from the user to a Device Token Flow
func PollForAccessToken(client http.Client, clientID string, deviceCodeResponse DeviceCodeResponse) (string, error) {
	expiry := time.Now().Add(time.Duration(deviceCodeResponse.ExpiresIn) * time.Second)
	interval := time.Duration(deviceCodeResponse.Interval) * time.Second

	var (
		accessToken string
		err         error
	)

	for time.Now().Before(expiry) {
		time.Sleep(interval)

		accessToken, err = requestAccessToken(client, clientID, deviceCodeResponse.DeviceCode)

		switch {
		case err == nil:
			return accessToken, nil
		case errors.Is(err, errPending):
			continue
		case errors.Is(err, errSlowDown):
			interval = interval + (5 * time.Second)
		case errors.Is(err, ErrExpiredToken):
			return "", fmt.Errorf("device token expired: %w", err)
		case errors.Is(err, ErrAccessDenied):
			return "", fmt.Errorf("user canceled: %w", err)
		default:
			return "", fmt.Errorf("requesting access token: %w", err)
		}
	}

	return "", errors.New("device token expired")
}
