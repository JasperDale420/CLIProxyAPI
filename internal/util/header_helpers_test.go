package util

import (
	"net/http"
	"testing"
)

func TestExtractCustomHeaders_FiltersAndTrims(t *testing.T) {
	got := extractCustomHeaders(map[string]string{
		"header:X-Keep":  "  v1 ",
		"header: X-Trim": "v2",
		"header:":        "novalue-name",
		"header:X-Empty": "   ",
		"api_key":        "not-a-header",
	})
	if len(got) != 2 || got["X-Keep"] != "v1" || got["X-Trim"] != "v2" {
		t.Fatalf("unexpected headers: %#v", got)
	}
}

func TestExtractCustomHeaders_NoMatchesReturnsNil(t *testing.T) {
	if got := extractCustomHeaders(map[string]string{"other": "x"}); got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if got := extractCustomHeaders(nil); got != nil {
		t.Fatalf("expected nil for nil attrs, got %#v", got)
	}
}

func TestApplyCustomHeadersFromAttrs_OverridesAndMirrorsHost(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("User-Agent", "default")
	ApplyCustomHeadersFromAttrs(r, map[string]string{
		"header:User-Agent": "custom",
		"header:Host":       "override.test",
	})
	if r.Header.Get("User-Agent") != "custom" {
		t.Fatalf("custom header should override default, got %q", r.Header.Get("User-Agent"))
	}
	if r.Host != "override.test" || r.Header.Get("Host") != "override.test" {
		t.Fatalf("Host not mirrored: Host=%q header=%q", r.Host, r.Header.Get("Host"))
	}
}

func TestApplyCustomHeadersFromAttrs_NilRequestIsSafe(t *testing.T) {
	ApplyCustomHeadersFromAttrs(nil, map[string]string{"header:X": "y"})
}
