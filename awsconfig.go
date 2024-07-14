package awsconsole

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
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

type ssoInfo struct {
	url         string
	sessionName string
}

func extractSSOInfo(filePath string, profileName string) (ssoInfo, error) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return ssoInfo{}, err
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
		return ssoInfo{}, fmt.Errorf("sso_start_url and sso_session for %q profile are not found", profileName)
	}

	return ssoInfo{
		url:         url,
		sessionName: sessionName,
	}, nil
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

func cachedTokenFilepath(info ssoInfo) (string, error) {
	if info.sessionName == "" {
		return ssocreds.StandardCachedTokenFilepath(info.url)
	} else {
		return ssocreds.StandardCachedTokenFilepath(info.sessionName)
	}
}
