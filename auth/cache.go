package auth

import (
	"encoding/json"
	"errors"

	"github.com/zalando/go-keyring"
)

var (
	keyringService = "dev.christopher2k.stream-chat-go"
	keyringUser    = "auth"
)

var (
	ErrorNoAuth = errors.New("No authentication")
)

func GetSavedAuth() (AuthTokenResponse, error) {
	savedAuth, err := keyring.Get(keyringService, keyringUser)
	if err != nil {
		return AuthTokenResponse{}, ErrorNoAuth
	}
	var auth AuthTokenResponse
	json.Unmarshal([]byte(savedAuth), &auth)
	return auth, nil
}

func SaveAuth(auth AuthTokenResponse) error {
	authJson, _ := json.Marshal(auth)
	return keyring.Set(keyringService, keyringUser, string(authJson))
}
