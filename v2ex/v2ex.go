// Package v2ex is the library behind the v2ex command: the HTTP client,
// request shaping, and the typed data models for V2EX.
//
// The client GETs resources from the public V2EX JSON API at
// https://www.v2ex.com/api. No authentication is required. It sets a
// real User-Agent, paces requests, and retries transient 429/5xx errors with
// exponential backoff.
package v2ex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// DefaultUserAgent identifies the client to V2EX.
const DefaultUserAgent = "v2ex/dev (+https://github.com/tamnd/v2ex-cli)"

// ErrNotFound is returned when the API returns no results for a single item.
var ErrNotFound = errors.New("not found")

// Config holds constructor parameters.
type Config struct {
	BaseURL   string
	UserAgent string
	Rate      time.Duration
	Retries   int
	Timeout   time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://www.v2ex.com",
		UserAgent: DefaultUserAgent,
		Rate:      500 * time.Millisecond,
		Retries:   3,
		Timeout:   30 * time.Second,
	}
}

// Client talks to the V2EX JSON API.
type Client struct {
	cfg        Config
	httpClient *http.Client
	mu         sync.Mutex
	last       time.Time
}

// NewClient returns a Client with the given config.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

// get fetches the given path (relative to BaseURL) and decodes JSON into dst.
func (c *Client) get(ctx context.Context, path string, dst any) error {
	fullURL := c.cfg.BaseURL + path

	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		retry, err := c.do(ctx, fullURL, dst)
		if err == nil {
			return nil
		}
		lastErr = err
		if !retry {
			return err
		}
	}
	return fmt.Errorf("get %s: %w", path, lastErr)
}

func (c *Client) do(ctx context.Context, rawURL string, dst any) (retry bool, err error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("http %d", resp.StatusCode)
	}

	b, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return true, err
	}
	if err := json.Unmarshal(b, dst); err != nil {
		return false, fmt.Errorf("decode: %w", err)
	}
	return false, nil
}

func (c *Client) pace() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 500 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}

// ─── public API ───────────────────────────────────────────────────────────────

// Topics fetches a topic list from the given API endpoint path
// (e.g. "/api/topics/hot.json" or "/api/topics/latest.json").
func (c *Client) Topics(ctx context.Context, endpoint string) ([]Topic, error) {
	var topics []Topic
	if err := c.get(ctx, endpoint, &topics); err != nil {
		return nil, err
	}
	return topics, nil
}

// Topic fetches a single topic by numeric ID string.
// Returns ErrNotFound if the API returns an empty array.
func (c *Client) Topic(ctx context.Context, id string) (*Topic, error) {
	path := "/api/topics/show.json?id=" + url.QueryEscape(id)
	var topics []Topic
	if err := c.get(ctx, path, &topics); err != nil {
		return nil, err
	}
	if len(topics) == 0 {
		return nil, ErrNotFound
	}
	return &topics[0], nil
}

// Node fetches a single node by name.
// Returns ErrNotFound when the API responds with HTTP 404.
func (c *Client) Node(ctx context.Context, name string) (*Node, error) {
	path := "/api/nodes/show.json?name=" + url.QueryEscape(name)
	var node Node
	if err := c.get(ctx, path, &node); err != nil {
		return nil, err
	}
	if node.ID == 0 {
		return nil, ErrNotFound
	}
	return &node, nil
}

// Member fetches a member profile by username.
// Returns ErrNotFound when the API responds with HTTP 404.
func (c *Client) Member(ctx context.Context, username string) (*Member, error) {
	path := "/api/members/show.json?username=" + url.QueryEscape(username)
	var member Member
	if err := c.get(ctx, path, &member); err != nil {
		return nil, err
	}
	if member.ID == 0 {
		return nil, ErrNotFound
	}
	return &member, nil
}

// Replies fetches all replies for a topic by topic ID string.
func (c *Client) Replies(ctx context.Context, topicID string) ([]Reply, error) {
	path := "/api/replies/show.json?topic_id=" + url.QueryEscape(topicID)
	var replies []Reply
	if err := c.get(ctx, path, &replies); err != nil {
		return nil, err
	}
	return replies, nil
}

// TopicsByNode fetches the most recent topics for a node by name.
// Returns ErrNotFound when the API returns an empty array.
func (c *Client) TopicsByNode(ctx context.Context, nodeName string) ([]Topic, error) {
	path := "/api/topics/show.json?node_name=" + url.QueryEscape(nodeName)
	var topics []Topic
	if err := c.get(ctx, path, &topics); err != nil {
		return nil, err
	}
	if len(topics) == 0 {
		return nil, ErrNotFound
	}
	return topics, nil
}

// AllNodes fetches the complete list of V2EX nodes.
func (c *Client) AllNodes(ctx context.Context) ([]Node, error) {
	var nodes []Node
	if err := c.get(ctx, "/api/nodes/all.json", &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

// ─── display helpers ──────────────────────────────────────────────────────────

// TopicToRow converts a Topic to a display-friendly TopicRow.
func TopicToRow(t Topic, rank int) TopicRow {
	return TopicRow{
		Rank:    rank,
		ID:      t.ID,
		Title:   t.Title,
		Node:    t.Node.Title,
		Replies: t.Replies,
		URL:     t.URL,
	}
}

// TopicToDetails returns field/value rows for a single topic view.
func TopicToDetails(t Topic) []TopicDetail {
	return []TopicDetail{
		{Field: "id", Value: fmt.Sprintf("%d", t.ID)},
		{Field: "title", Value: t.Title},
		{Field: "author", Value: t.Member.Username},
		{Field: "node", Value: t.Node.Title},
		{Field: "replies", Value: fmt.Sprintf("%d", t.Replies)},
		{Field: "created", Value: unixDate(t.Created)},
		{Field: "content", Value: t.Content},
		{Field: "url", Value: t.URL},
	}
}

// NodeToDetails returns field/value rows for a node view.
func NodeToDetails(n Node) []NodeDetail {
	return []NodeDetail{
		{Field: "id", Value: fmt.Sprintf("%d", n.ID)},
		{Field: "name", Value: n.Name},
		{Field: "title", Value: n.Title},
		{Field: "topics", Value: fmt.Sprintf("%d", n.Topics)},
		{Field: "stars", Value: fmt.Sprintf("%d", n.Stars)},
	}
}

// MemberToDetails returns field/value rows for a member view.
func MemberToDetails(m Member) []MemberDetail {
	tagline := m.Tagline
	if tagline == "" {
		tagline = m.Bio
	}
	return []MemberDetail{
		{Field: "username", Value: m.Username},
		{Field: "tagline", Value: tagline},
		{Field: "created", Value: unixDate(m.Created)},
	}
}

// NodeToRow converts a Node to the flat NodeRow used in all-nodes lists.
func NodeToRow(n Node, rank int) NodeRow {
	return NodeRow{
		Rank:   rank,
		ID:     n.ID,
		Name:   n.Name,
		Title:  n.Title,
		Topics: n.Topics,
		Stars:  n.Stars,
	}
}

// ReplyToRow converts a Reply to a display-friendly ReplyRow.
func ReplyToRow(r Reply, rank int) ReplyRow {
	return ReplyRow{
		Rank:      rank,
		ID:        r.ID,
		Author:    r.Member.Username,
		CreatedAt: unixDate(r.Created),
		Content:   r.Content,
	}
}
