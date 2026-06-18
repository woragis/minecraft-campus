package campusworld

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	serverSlug string
	http       *http.Client
}

type WhitelistResult struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
}

func NewClient(baseURL, apiKey, serverSlug string) *Client {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if serverSlug == "" {
		serverSlug = "bedrock"
	}
	return &Client{
		baseURL:    baseURL,
		apiKey:     strings.TrimSpace(apiKey),
		serverSlug: serverSlug,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) CheckBedrockWhitelist(ctx context.Context, xuid, username string) (*WhitelistResult, error) {
	xuid = strings.TrimSpace(xuid)
	reqURL := fmt.Sprintf("%s/v1/internal/whitelist/bedrock/%s?username=%s", c.baseURL, url.PathEscape(xuid), url.QueryEscape(username))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("X-Plugin-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("whitelist request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read whitelist response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("whitelist status %d: %s", resp.StatusCode, string(body))
	}

	var out WhitelistResult
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode whitelist: %w", err)
	}
	return &out, nil
}

func (c *Client) UpsertBedrockPlayer(ctx context.Context, xuid, username string) error {
	payload := map[string]string{
		"xuid":       strings.TrimSpace(xuid),
		"username":   strings.TrimSpace(username),
		"serverSlug": c.serverSlug,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode upsert: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/internal/players/bedrock/upsert", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("new upsert request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Plugin-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("upsert request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read upsert response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upsert status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func KickMessage(reason string) string {
	switch reason {
	case "not_invited":
		return "You need a CampusWorld invite to join. Ask a member for an invite."
	case "banned":
		return "You are banned from CampusWorld."
	default:
		return "You cannot join CampusWorld right now."
	}
}
