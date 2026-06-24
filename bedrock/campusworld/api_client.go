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

type PlayerUpsertResult struct {
	ID       string `json:"id"`
	Username string `json:"username"`
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

func (c *Client) UpsertBedrockPlayer(ctx context.Context, xuid, username string) (*PlayerUpsertResult, error) {
	payload := map[string]string{
		"xuid":       strings.TrimSpace(xuid),
		"username":   strings.TrimSpace(username),
		"serverSlug": c.serverSlug,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode upsert: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/internal/players/bedrock/upsert", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("new upsert request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Plugin-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upsert request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read upsert response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upsert status %d: %s", resp.StatusCode, string(body))
	}
	var out PlayerUpsertResult
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode upsert: %w", err)
	}
	return &out, nil
}

func (c *Client) PresenceOffline(ctx context.Context, playerID string) error {
	return c.postPresence(ctx, "/v1/internal/presence/offline", playerID)
}

func (c *Client) PresenceHeartbeat(ctx context.Context, playerID string) error {
	return c.postPresence(ctx, "/v1/internal/presence/heartbeat", playerID)
}

type HUDResult struct {
	Username         string `json:"username"`
	Status           string `json:"status"`
	AffiliationType  string `json:"affiliationType"`
	UniversitySlug   string `json:"universitySlug,omitempty"`
	FacultySlug      string `json:"facultySlug,omitempty"`
	CourseSlug       string `json:"courseSlug,omitempty"`
	UniversityName   string `json:"universityName,omitempty"`
	UniversityHex    string `json:"universityHex,omitempty"`
	FacultyName      string `json:"facultyName,omitempty"`
	FacultyAbbr      string `json:"facultyAbbr,omitempty"`
	FacultyHex       string `json:"facultyHex,omitempty"`
	CourseName       string `json:"courseName,omitempty"`
	CourseAbbr       string `json:"courseAbbr,omitempty"`
	CourseHex        string `json:"courseHex,omitempty"`
	GuildID          string `json:"guildId,omitempty"`
	GuildName        string `json:"guildName,omitempty"`
	GuildSlug        string `json:"guildSlug,omitempty"`
	GuildOnlineCount int    `json:"guildOnlineCount"`
}

func (c *Client) StatsIngest(ctx context.Context, playerID string, sessionSeconds, mobKills int64) error {
	payload := map[string]any{
		"playerId":       strings.TrimSpace(playerID),
		"serverSlug":     c.serverSlug,
		"sessionSeconds": sessionSeconds,
		"mobKills":       mobKills,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode stats: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/internal/stats/ingest", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("new stats request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Plugin-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("stats request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read stats response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("stats status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *Client) FetchHUD(ctx context.Context, playerID string) (*HUDResult, error) {
	playerID = strings.TrimSpace(playerID)
	reqURL := fmt.Sprintf("%s/v1/internal/players/%s/hud", c.baseURL, url.PathEscape(playerID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new hud request: %w", err)
	}
	req.Header.Set("X-Plugin-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hud request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read hud response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hud status %d: %s", resp.StatusCode, string(body))
	}
	var out HUDResult
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode hud: %w", err)
	}
	return &out, nil
}

func (c *Client) postPresence(ctx context.Context, path, playerID string) error {
	payload := map[string]string{
		"playerId":   strings.TrimSpace(playerID),
		"serverSlug": c.serverSlug,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode presence: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("new presence request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Plugin-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("presence request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read presence response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("presence status %d: %s", resp.StatusCode, string(body))
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
