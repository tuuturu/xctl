package github

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// https://docs.github.com/en/developers/apps/building-oauth-apps/authorizing-oauth-apps#error-codes-for-the-device-flow
var (
	// errPending indicates user has yet to complete the flow
	errPending = errors.New("authorization pending")
	// errSlowDown indicates request interval must be throttled
	errSlowDown = errors.New("slow down")
	// ErrExpiredToken indicates device code expiry
	ErrExpiredToken = errors.New("expired token")
	// ErrAccessDenied indicates rejection by user
	ErrAccessDenied = errors.New("authorization rejected")
)

func asError(payload []byte) error {
	var potentialError struct {
		Error string `json:"error"`
	}

	err := json.Unmarshal(payload, &potentialError)
	if err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}

	if potentialError.Error == "" {
		return nil
	}

	switch potentialError.Error {
	case "authorization_pending":
		return errPending
	case "slow_down":
		return errSlowDown
	case "expired_token":
		return ErrExpiredToken
	case "access_denied":
		return ErrAccessDenied
	default:
		return err
	}
}
