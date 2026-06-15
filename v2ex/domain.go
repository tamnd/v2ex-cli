package v2ex

import (
	"context"
	"errors"
	"strconv"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

// domain.go exposes v2ex as a kit Domain. A multi-domain host (ant) enables it
// with a single blank import:
//
//	import _ "github.com/tamnd/v2ex-cli/v2ex"
//
// The same Domain also drives the standalone v2ex binary (see cli/) so the
// binary and a host share one source of truth.

func init() { kit.Register(Domain{}) }

// Domain is the v2ex driver. It carries no state; the per-run client is
// built by the factory Register installs via app.SetClient.
type Domain struct{}

// Info returns the canonical scheme, hostnames, and identity for v2ex.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme: "v2ex",
		Hosts:  []string{"www.v2ex.com", "v2ex.com"},
		Identity: kit.Identity{
			Binary: "v2ex",
			Short:  "A command line for V2EX.",
			Long: `A command line for V2EX (v2ex.com).

Browse hot and latest topics, look up nodes, member profiles, and replies
over the public V2EX JSON API. No API key required.`,
			Site: "https://www.v2ex.com",
			Repo: "https://github.com/tamnd/v2ex-cli",
		},
	}
}

// Register installs the client factory and all 6 operations onto app.
func (Domain) Register(app *kit.App) {
	app.SetClient(newClient)

	// hot: list hot topics
	kit.Handle(app, kit.OpMeta{
		Name:    "hot",
		Group:   "read",
		List:    true,
		Summary: "List hot topics",
		URIType: "topic",
	}, listHot)

	// latest: list latest topics
	kit.Handle(app, kit.OpMeta{
		Name:    "latest",
		Group:   "read",
		List:    true,
		Summary: "List latest topics",
		URIType: "topic",
	}, listLatest)

	// topic: fetch a single topic by ID
	kit.Handle(app, kit.OpMeta{
		Name:     "topic",
		Group:    "read",
		Single:   true,
		Resolver: true,
		Summary:  "Get a topic by ID",
		URIType:  "topic",
		Args:     []kit.Arg{{Name: "id", Help: "numeric topic ID"}},
	}, getTopic)

	// node: fetch a single node by name
	kit.Handle(app, kit.OpMeta{
		Name:     "node",
		Group:    "read",
		Single:   true,
		Resolver: true,
		Summary:  "Get a node by name",
		URIType:  "node",
		Args:     []kit.Arg{{Name: "name", Help: "node name (slug)"}},
	}, getNode)

	// member: fetch a member profile by username
	kit.Handle(app, kit.OpMeta{
		Name:     "member",
		Group:    "read",
		Single:   true,
		Resolver: true,
		Summary:  "Get a member profile by username",
		URIType:  "member",
		Args:     []kit.Arg{{Name: "username", Help: "member username"}},
	}, getMember)

	// replies: list replies for a topic
	kit.Handle(app, kit.OpMeta{
		Name:    "replies",
		Group:   "read",
		List:    true,
		Summary: "List replies for a topic",
		URIType: "topic",
		Args:    []kit.Arg{{Name: "topic_id", Help: "numeric topic ID"}},
	}, listReplies)
}

// newClient builds a *Client from the kit-resolved config.
func newClient(_ context.Context, cfg kit.Config) (any, error) {
	c := Config{
		BaseURL:   "https://www.v2ex.com",
		UserAgent: DefaultUserAgent,
		Rate:      DefaultConfig().Rate,
		Retries:   DefaultConfig().Retries,
		Timeout:   DefaultConfig().Timeout,
	}
	if cfg.UserAgent != "" {
		c.UserAgent = cfg.UserAgent
	}
	if cfg.Rate > 0 {
		c.Rate = cfg.Rate
	}
	if cfg.Retries > 0 {
		c.Retries = cfg.Retries
	}
	if cfg.Timeout > 0 {
		c.Timeout = cfg.Timeout
	}
	return NewClient(c), nil
}

// --- input structs ---

type noArgs struct {
	Client *Client `kit:"inject"`
	Limit  int     `kit:"flag,inherit" help:"max results"`
}

type topicArgs struct {
	ID     string  `kit:"arg" help:"numeric topic ID"`
	Client *Client `kit:"inject"`
}

type nodeArgs struct {
	Name   string  `kit:"arg" help:"node name (slug)"`
	Client *Client `kit:"inject"`
}

type memberArgs struct {
	Username string  `kit:"arg" help:"member username"`
	Client   *Client `kit:"inject"`
}

type repliesArgs struct {
	TopicID string  `kit:"arg" help:"numeric topic ID"`
	Client  *Client `kit:"inject"`
	Limit   int     `kit:"flag,inherit" help:"max results"`
}

// --- handlers ---

func listHot(ctx context.Context, in noArgs, emit func(*Topic) error) error {
	topics, err := in.Client.Topics(ctx, "/api/topics/hot.json")
	if err != nil {
		return mapErr(err)
	}
	for i := range topics {
		if in.Limit > 0 && i >= in.Limit {
			break
		}
		if err := emit(&topics[i]); err != nil {
			return err
		}
	}
	return nil
}

func listLatest(ctx context.Context, in noArgs, emit func(*Topic) error) error {
	topics, err := in.Client.Topics(ctx, "/api/topics/latest.json")
	if err != nil {
		return mapErr(err)
	}
	for i := range topics {
		if in.Limit > 0 && i >= in.Limit {
			break
		}
		if err := emit(&topics[i]); err != nil {
			return err
		}
	}
	return nil
}

func getTopic(ctx context.Context, in topicArgs, emit func(*Topic) error) error {
	t, err := in.Client.Topic(ctx, in.ID)
	if err != nil {
		return mapErr(err)
	}
	return emit(t)
}

func getNode(ctx context.Context, in nodeArgs, emit func(*Node) error) error {
	n, err := in.Client.Node(ctx, in.Name)
	if err != nil {
		return mapErr(err)
	}
	return emit(n)
}

func getMember(ctx context.Context, in memberArgs, emit func(*Member) error) error {
	m, err := in.Client.Member(ctx, in.Username)
	if err != nil {
		return mapErr(err)
	}
	return emit(m)
}

func listReplies(ctx context.Context, in repliesArgs, emit func(*Reply) error) error {
	replies, err := in.Client.Replies(ctx, in.TopicID)
	if err != nil {
		return mapErr(err)
	}
	for i := range replies {
		if in.Limit > 0 && i >= in.Limit {
			break
		}
		if err := emit(&replies[i]); err != nil {
			return err
		}
	}
	return nil
}

// --- Resolver: pure string URI functions ---

// Classify turns a bare id or https URL into (uriType, id).
func (Domain) Classify(input string) (uriType, id string, err error) {
	if _, parseErr := strconv.Atoi(input); parseErr == nil {
		return "topic", input, nil
	}
	return "", "", errs.Usage("unrecognized v2ex reference: %q (want a numeric topic ID)", input)
}

// Locate returns the https URL for a (uriType, id).
func (Domain) Locate(uriType, id string) (string, error) {
	switch uriType {
	case "topic":
		return "https://www.v2ex.com/t/" + id, nil
	case "node":
		return "https://www.v2ex.com/go/" + id, nil
	case "member":
		return "https://www.v2ex.com/member/" + id, nil
	default:
		return "", errs.Usage("v2ex has no resource type %q", uriType)
	}
}

// mapErr converts library errors to kit error kinds with the right exit codes.
func mapErr(err error) error {
	if errors.Is(err, ErrNotFound) {
		return errs.NotFound("%s", err.Error())
	}
	return err
}
