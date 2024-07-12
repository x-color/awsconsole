package awsconsole

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

type token struct {
	AccessToken string    `json:"accessToken,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
}

func extractAccessToken(filePath string) (string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	tk := token{}
	if err := json.Unmarshal(b, &tk); err != nil {
		return "", err
	}

	if time.Now().After(tk.ExpiresAt) {
		return "", fmt.Errorf("token expired")
	}

	return tk.AccessToken, nil
}

func extractSSOStartURL(filePath string, profileName string) (string, error) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return "", err
	}

	sectionName := "default"
	if profileName != "default" {
		sectionName = "profile " + profileName
	}

	profile := cfg.Section(sectionName)
	sessionName := profile.Key("sso_session").String()

	url := cfg.Section("sso-session " + sessionName).Key("sso_start_url").String()
	if url == "" {
		url = profile.Key("sso_start_url").String()
	}

	if url == "" {
		return "", fmt.Errorf("sso_start_url and sso_session for %q profile are not found", profileName)
	}

	return url, nil
}