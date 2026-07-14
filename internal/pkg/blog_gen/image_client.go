package blog_gen

import (
	"bytes"
	"deller/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"deller/internal/service"
)

type imageChatRequest struct {
	Model        string        `json:"model"`
	Messages     []chatMessage `json:"messages"`
	FunctionCall string        `json:"function_call"`
}

func fetchImageBytes(prompt string) ([]byte, error) {
	token, err := getOAuthToken()
	if err != nil {
		return nil, fmt.Errorf("get oauth token: %w", err)
	}

	reqBody := imageChatRequest{
		Model: config.GetGigaChatModel(),
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
		FunctionCall: "auto",
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal image request: %w", err)
	}

	apiURL := config.GetGigaChatAPIURL() + "/chat/completions"
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create image request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := newGigaChatClient(180 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send image request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read image response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		clearTokenCache()
		token, err = getOAuthToken()
		if err != nil {
			return nil, fmt.Errorf("refresh oauth token: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("retry image request: %w", err)
		}
		defer resp.Body.Close()

		respBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read retry image response: %w", err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("image API returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal image response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("image API returned no choices")
	}

	content := chatResp.Choices[0].Message.Content
	fileID := extractFileID(content)
	if fileID == "" {
		return nil, fmt.Errorf("no file_id found in response: %s", content)
	}

	return downloadFile(fileID, token)
}

func extractFileID(content string) string {
	re := regexp.MustCompile(`src="([^"]+)"`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}

	re2 := regexp.MustCompile(`([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)
	matches2 := re2.FindStringSubmatch(content)
	if len(matches2) > 1 {
		return matches2[1]
	}

	return ""
}

func downloadFile(fileID, token string) ([]byte, error) {
	apiURL := config.GetGigaChatAPIURL() + "/files/" + fileID + "/content"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create file request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := newGigaChatClient(60 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send file request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		clearTokenCache()
		token, err = getOAuthToken()
		if err != nil {
			return nil, fmt.Errorf("refresh oauth token: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("retry file request: %w", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("file API returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	return io.ReadAll(resp.Body)
}

func saveImage(imgBytes []byte) (string, error) {
	imagesDir := filepath.Join("shared", "images")
	if err := os.MkdirAll(imagesDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("create images dir: %w", err)
	}

	uniqueName, err := service.GenerateUniqueFileName("blog_post.png", imagesDir)
	if err != nil {
		return "", fmt.Errorf("generate unique filename: %w", err)
	}

	filePath := filepath.Join(imagesDir, uniqueName)
	if err := os.WriteFile(filePath, imgBytes, 0644); err != nil {
		return "", fmt.Errorf("save image: %w", err)
	}

	return "/images/" + uniqueName, nil
}
