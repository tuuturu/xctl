package github

import (
	"github.com/google/go-github/v44/github"
)

type authenticationService struct {
	client *github.Client
}

const (
	grantType = "urn:ietf:params:oauth:grant-type:device_code"
)

type accessTokenRequest struct {
	ClientID   string `json:"client_id"`
	GrantType  string `json:"grant_type"`
	DeviceCode string `json:"device_code"`
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type deviceCodeRequest struct {
	ClientID string `json:"client_id"`
	Scope    string `json:"scope"`
}
