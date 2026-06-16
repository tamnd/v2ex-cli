package v2ex

import "time"

// Topic is a V2EX discussion thread.
type Topic struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	URL     string `json:"url"`
	Replies int    `json:"replies"`
	Content string `json:"content"`
	Created int64  `json:"created"`
	Node    struct {
		Title string `json:"title"`
	} `json:"node"`
	Member struct {
		Username string `json:"username"`
	} `json:"member"`
}

// Node is a V2EX topic category.
type Node struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Topics int    `json:"topics"`
	Stars  int    `json:"stars"`
}

// Member is a V2EX registered user.
type Member struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Bio          string `json:"bio"`
	Tagline      string `json:"tagline"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"last_modified"`
}

// Reply is a comment posted under a topic.
type Reply struct {
	ID      int   `json:"id"`
	Content string `json:"content"`
	Created int64  `json:"created"`
	Member  struct {
		Username string `json:"username"`
	} `json:"member"`
}

// ─── display records ──────────────────────────────────────────────────────────

// TopicRow is the flat record rendered for hot/latest lists.
type TopicRow struct {
	Rank    int    `json:"rank"`
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Node    string `json:"node"`
	Replies int    `json:"replies"`
	URL     string `json:"url"`
}

// TopicDetail is the flat record rendered for a single topic view.
type TopicDetail struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

// NodeDetail is a flat field/value record for a node.
type NodeDetail = TopicDetail

// MemberDetail is a flat field/value record for a member.
type MemberDetail = TopicDetail

// ReplyRow is the flat record rendered for the replies list.
type ReplyRow struct {
	Rank      int    `json:"rank"`
	ID        int    `json:"id"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
	Content   string `json:"content"`
}

// NodeRow is the flat record rendered for the all-nodes list.
type NodeRow struct {
	Rank   int    `json:"rank"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Topics int    `json:"topics"`
	Stars  int    `json:"stars"`
}

// unixDate converts a Unix timestamp to YYYY-MM-DD.
func unixDate(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).UTC().Format("2006-01-02")
}
