package v2ex

import (
	"testing"

	"github.com/tamnd/any-cli/kit"
)

func TestDomainInfo(t *testing.T) {
	info := Domain{}.Info()
	if info.Scheme != "v2ex" {
		t.Errorf("Scheme = %q, want v2ex", info.Scheme)
	}
	if info.Identity.Binary != "v2ex" {
		t.Errorf("Binary = %q, want v2ex", info.Identity.Binary)
	}
}

func TestDomainRegistered(t *testing.T) {
	h, err := kit.Open()
	if err != nil {
		t.Fatal(err)
	}
	domains := h.Domains()
	for _, d := range domains {
		if d == "v2ex" {
			return
		}
	}
	t.Error("v2ex domain not registered")
}

func TestLocate(t *testing.T) {
	cases := []struct {
		uriType string
		id      string
		want    string
	}{
		{"topic", "123", "https://www.v2ex.com/t/123"},
		{"node", "golang", "https://www.v2ex.com/go/golang"},
		{"member", "alice", "https://www.v2ex.com/member/alice"},
	}
	for _, tc := range cases {
		got, err := Domain{}.Locate(tc.uriType, tc.id)
		if err != nil || got != tc.want {
			t.Errorf("Locate(%q, %q) = (%q, %v), want (%q, nil)", tc.uriType, tc.id, got, err, tc.want)
		}
	}
}

func TestClassify(t *testing.T) {
	typ, id, err := Domain{}.Classify("12345")
	if err != nil || typ != "topic" || id != "12345" {
		t.Errorf("Classify(12345) = (%q, %q, %v), want (topic, 12345, nil)", typ, id, err)
	}

	_, _, err = Domain{}.Classify("not-a-number")
	if err == nil {
		t.Error("Classify(not-a-number) should return error")
	}
}
