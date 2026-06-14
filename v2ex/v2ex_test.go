package v2ex_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/v2ex-cli/v2ex"
)

func newTestClient(ts *httptest.Server) *v2ex.Client {
	cfg := v2ex.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return v2ex.NewClient(cfg)
}

// ─── mock responses ───────────────────────────────────────────────────────────

const mockTopicsJSON = `[
  {
    "id": 1001,
    "title": "Welcome to V2EX",
    "url": "https://www.v2ex.com/t/1001",
    "replies": 10,
    "content": "Hello world",
    "created": 1700000000,
    "node": {"title": "Go"},
    "member": {"username": "alice"}
  },
  {
    "id": 1002,
    "title": "Another topic",
    "url": "https://www.v2ex.com/t/1002",
    "replies": 5,
    "content": "",
    "created": 1700001000,
    "node": {"title": "Python"},
    "member": {"username": "bob"}
  }
]`

const mockTopicJSON = `[
  {
    "id": 1001,
    "title": "Welcome to V2EX",
    "url": "https://www.v2ex.com/t/1001",
    "replies": 10,
    "content": "Hello world",
    "created": 1700000000,
    "node": {"title": "Go"},
    "member": {"username": "alice"}
  }
]`

const mockEmptyTopicJSON = `[]`

const mockNodeJSON = `{
  "id": 42,
  "name": "go",
  "title": "Go Programming",
  "topics": 1500,
  "stars": 8000
}`

const mockMemberJSON = `{
  "id": 99,
  "username": "alice",
  "bio": "Developer",
  "tagline": "I write Go",
  "created": 1300000000,
  "last_modified": 1700000000
}`

const mockMemberNoBioJSON = `{
  "id": 99,
  "username": "bob",
  "bio": "just bio",
  "tagline": "",
  "created": 1300000000,
  "last_modified": 1700000000
}`

const mockRepliesJSON = `[
  {
    "id": 501,
    "content": "Great post!",
    "created": 1700000100,
    "member": {"username": "charlie"}
  },
  {
    "id": 502,
    "content": "Thanks for sharing",
    "created": 1700000200,
    "member": {"username": "dave"}
  }
]`

// ─── Topics ───────────────────────────────────────────────────────────────────

func TestTopicsHot(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/topics/hot.json" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(mockTopicsJSON))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	topics, err := c.Topics(context.Background(), "/api/topics/hot.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(topics) != 2 {
		t.Fatalf("got %d topics, want 2", len(topics))
	}
	if topics[0].ID != 1001 {
		t.Errorf("id = %d, want 1001", topics[0].ID)
	}
	if topics[0].Title != "Welcome to V2EX" {
		t.Errorf("title = %q", topics[0].Title)
	}
	if topics[0].Node.Title != "Go" {
		t.Errorf("node = %q, want Go", topics[0].Node.Title)
	}
	if topics[0].Member.Username != "alice" {
		t.Errorf("member = %q, want alice", topics[0].Member.Username)
	}
}

func TestTopicsLatest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/topics/latest.json" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(mockTopicsJSON))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	topics, err := c.Topics(context.Background(), "/api/topics/latest.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(topics) == 0 {
		t.Fatal("got 0 topics")
	}
}

// ─── Topic ────────────────────────────────────────────────────────────────────

func TestTopicShow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/topics/show.json" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if r.URL.Query().Get("id") != "1001" {
			t.Errorf("id query = %q, want 1001", r.URL.Query().Get("id"))
		}
		_, _ = w.Write([]byte(mockTopicJSON))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	topic, err := c.Topic(context.Background(), "1001")
	if err != nil {
		t.Fatal(err)
	}
	if topic.ID != 1001 {
		t.Errorf("id = %d, want 1001", topic.ID)
	}
	if topic.Title != "Welcome to V2EX" {
		t.Errorf("title = %q", topic.Title)
	}
}

func TestTopicNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(mockEmptyTopicJSON))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.Topic(context.Background(), "9999")
	if !errors.Is(err, v2ex.ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

// ─── Node ─────────────────────────────────────────────────────────────────────

func TestNodeShow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/nodes/show.json" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if r.URL.Query().Get("name") != "go" {
			t.Errorf("name query = %q, want go", r.URL.Query().Get("name"))
		}
		_, _ = w.Write([]byte(mockNodeJSON))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	node, err := c.Node(context.Background(), "go")
	if err != nil {
		t.Fatal(err)
	}
	if node.ID != 42 {
		t.Errorf("id = %d, want 42", node.ID)
	}
	if node.Name != "go" {
		t.Errorf("name = %q, want go", node.Name)
	}
	if node.Topics != 1500 {
		t.Errorf("topics = %d, want 1500", node.Topics)
	}
	if node.Stars != 8000 {
		t.Errorf("stars = %d, want 8000", node.Stars)
	}
}

// ─── Member ───────────────────────────────────────────────────────────────────

func TestMemberShow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/members/show.json" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if r.URL.Query().Get("username") != "alice" {
			t.Errorf("username query = %q, want alice", r.URL.Query().Get("username"))
		}
		_, _ = w.Write([]byte(mockMemberJSON))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	member, err := c.Member(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if member.ID != 99 {
		t.Errorf("id = %d, want 99", member.ID)
	}
	if member.Username != "alice" {
		t.Errorf("username = %q, want alice", member.Username)
	}
	if member.Tagline != "I write Go" {
		t.Errorf("tagline = %q", member.Tagline)
	}
}

func TestMemberTaglineFallbackToBio(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(mockMemberNoBioJSON))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	member, err := c.Member(context.Background(), "bob")
	if err != nil {
		t.Fatal(err)
	}
	details := v2ex.MemberToDetails(*member)
	var taglineVal string
	for _, d := range details {
		if d.Field == "tagline" {
			taglineVal = d.Value
		}
	}
	if taglineVal != "just bio" {
		t.Errorf("tagline fallback = %q, want 'just bio'", taglineVal)
	}
}

// ─── Replies ──────────────────────────────────────────────────────────────────

func TestRepliesShow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/replies/show.json" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if r.URL.Query().Get("topic_id") != "1001" {
			t.Errorf("topic_id query = %q, want 1001", r.URL.Query().Get("topic_id"))
		}
		_, _ = w.Write([]byte(mockRepliesJSON))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	replies, err := c.Replies(context.Background(), "1001")
	if err != nil {
		t.Fatal(err)
	}
	if len(replies) != 2 {
		t.Fatalf("got %d replies, want 2", len(replies))
	}
	if replies[0].ID != 501 {
		t.Errorf("id = %d, want 501", replies[0].ID)
	}
	if replies[0].Member.Username != "charlie" {
		t.Errorf("member = %q, want charlie", replies[0].Member.Username)
	}
}

// ─── retry ────────────────────────────────────────────────────────────────────

func TestClientRetries503(t *testing.T) {
	var hits int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if hits < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(mockTopicsJSON))
	}))
	defer ts.Close()

	cfg := v2ex.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	cfg.Retries = 5
	c := v2ex.NewClient(cfg)

	_, err := c.Topics(context.Background(), "/api/topics/hot.json")
	if err != nil {
		t.Fatal(err)
	}
	if hits != 3 {
		t.Errorf("server saw %d hits, want 3", hits)
	}
}

// ─── display helpers ──────────────────────────────────────────────────────────

func TestTopicToRow(t *testing.T) {
	topic := v2ex.Topic{
		ID:      1001,
		Title:   "Hello",
		URL:     "https://www.v2ex.com/t/1001",
		Replies: 7,
	}
	topic.Node.Title = "Go"

	row := v2ex.TopicToRow(topic, 3)
	if row.Rank != 3 {
		t.Errorf("rank = %d, want 3", row.Rank)
	}
	if row.Node != "Go" {
		t.Errorf("node = %q, want Go", row.Node)
	}
}

func TestReplyToRow(t *testing.T) {
	reply := v2ex.Reply{
		ID:      501,
		Content: "Nice",
		Created: 1700000100,
	}
	reply.Member.Username = "charlie"

	row := v2ex.ReplyToRow(reply, 1)
	if row.Author != "charlie" {
		t.Errorf("author = %q, want charlie", row.Author)
	}
	if row.CreatedAt == "" {
		t.Error("created_at should not be empty")
	}
}
