package blog_gen

import (
	"bytes"
	"crypto/tls"
	"deller/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

type oauthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

var (
	cachedToken     string
	cachedExpiresAt int64
	tokenMutex      sync.Mutex
)

func clearTokenCache() {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	fmt.Println("[BLOG_GEN] clearing token cache")
	cachedToken = ""
	cachedExpiresAt = 0
}

func getOAuthToken() (string, error) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	if cachedToken != "" && time.Now().Unix() < cachedExpiresAt-60 {
		return cachedToken, nil
	}

	apiKey := config.GetGigaChatAPIKey()
	scope := config.GetGigaChatScope()
	oauthURL := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"

	formData := url.Values{}
	formData.Set("scope", scope)

	req, err := http.NewRequest("POST", oauthURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("create oauth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+apiKey)
	rqUID, _ := uuid.NewV7()
	req.Header.Set("RqUID", rqUID.String())

	client := newGigaChatClient(30 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send oauth request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read oauth response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("oauth returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var oauthResp oauthResponse
	if err := json.Unmarshal(respBytes, &oauthResp); err != nil {
		return "", fmt.Errorf("unmarshal oauth response: %w", err)
	}

	cachedToken = oauthResp.AccessToken
	cachedExpiresAt = oauthResp.ExpiresAt

	fmt.Printf("[BLOG_GEN] token obtained, expires_at: %d, now: %d\n", cachedExpiresAt, time.Now().Unix())

	return cachedToken, nil
}

func newGigaChatClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func callLLM(systemPrompt, userPrompt string) (string, error) {
	token, err := getOAuthToken()
	if err != nil {
		return "", fmt.Errorf("get oauth token: %w", err)
	}

	reqBody := chatRequest{
		Model: config.GetGigaChatModel(),
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	apiURL := config.GetGigaChatAPIURL() + "/chat/completions"
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := newGigaChatClient(120 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		clearTokenCache()
		token, err = getOAuthToken()
		if err != nil {
			return "", fmt.Errorf("refresh oauth token: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		resp, err = client.Do(req)
		if err != nil {
			return "", fmt.Errorf("retry request: %w", err)
		}
		defer resp.Body.Close()

		respBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("read retry response: %w", err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("LLM returned no choices")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func extractJSON(raw string) (string, error) {
	start := -1
	end := -1
	for i, ch := range raw {
		if ch == '{' && start == -1 {
			start = i
		}
		if ch == '}' {
			end = i
		}
	}
	if start == -1 || end == -1 {
		return "", fmt.Errorf("no JSON object found in response")
	}
	return raw[start : end+1], nil
}
