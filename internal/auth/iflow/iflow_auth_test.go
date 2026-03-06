package iflow

import "testing"

func TestDefaultAPIBaseURLUsesCurrentZAIHost(t *testing.T) {
	const want = "https://api.z.ai/api/paas/v4"
	if DefaultAPIBaseURL != want {
		t.Fatalf("DefaultAPIBaseURL = %q, want %q", DefaultAPIBaseURL, want)
	}
}
