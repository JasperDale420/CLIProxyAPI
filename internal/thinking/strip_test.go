package thinking

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestStripThinkingConfig_RemovesProviderSpecificFields(t *testing.T) {
	tests := []struct {
		name      string
		provider  string
		body      string
		checkGone []string
		checkKept []string
	}{
		{
			name:      "claude removes thinking and effort, drops empty output_config",
			provider:  "claude",
			body:      `{"model":"claude-sonnet-4-5","thinking":{"type":"enabled","budget_tokens":8192},"output_config":{"effort":"high"}}`,
			checkGone: []string{"thinking", "output_config"},
			checkKept: []string{"model"},
		},
		{
			name:      "claude keeps output_config when other fields remain",
			provider:  "claude",
			body:      `{"output_config":{"effort":"high","other":"keep"}}`,
			checkGone: []string{"output_config.effort", "thinking"},
			checkKept: []string{"output_config.other"},
		},
		{
			name:      "gemini removes generationConfig.thinkingConfig only",
			provider:  "gemini",
			body:      `{"generationConfig":{"thinkingConfig":{"thinkingBudget":100},"temperature":0.5}}`,
			checkGone: []string{"generationConfig.thinkingConfig"},
			checkKept: []string{"generationConfig.temperature"},
		},
		{
			name:      "gemini-cli removes nested request path",
			provider:  "gemini-cli",
			body:      `{"request":{"generationConfig":{"thinkingConfig":{"thinkingBudget":100}}}}`,
			checkGone: []string{"request.generationConfig.thinkingConfig"},
		},
		{
			name:      "openai removes reasoning_effort",
			provider:  "openai",
			body:      `{"reasoning_effort":"high","model":"gpt-5.4"}`,
			checkGone: []string{"reasoning_effort"},
			checkKept: []string{"model"},
		},
		{
			name:      "kimi removes both reasoning_effort and thinking",
			provider:  "kimi",
			body:      `{"reasoning_effort":"high","thinking":{"type":"enabled"}}`,
			checkGone: []string{"reasoning_effort", "thinking"},
		},
		{
			name:      "codex removes reasoning.effort",
			provider:  "codex",
			body:      `{"reasoning":{"effort":"high","summary":"auto"}}`,
			checkGone: []string{"reasoning.effort"},
			checkKept: []string{"reasoning.summary"},
		},
		{
			name:      "unknown provider leaves body unchanged",
			provider:  "unknown-provider",
			body:      `{"thinking":{"type":"enabled"}}`,
			checkKept: []string{"thinking"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := StripThinkingConfig([]byte(tt.body), tt.provider)
			for _, path := range tt.checkGone {
				if gjson.GetBytes(out, path).Exists() {
					t.Errorf("expected %q to be removed, body=%s", path, string(out))
				}
			}
			for _, path := range tt.checkKept {
				if !gjson.GetBytes(out, path).Exists() {
					t.Errorf("expected %q to be kept, body=%s", path, string(out))
				}
			}
		})
	}
}

func TestStripThinkingConfig_HandlesInvalidInputUnchanged(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		body     []byte
	}{
		{"empty body", "claude", []byte{}},
		{"invalid json", "claude", []byte(`not json`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := StripThinkingConfig(tt.body, tt.provider)
			if string(out) != string(tt.body) {
				t.Errorf("expected body unchanged for %s, got %q", tt.name, string(out))
			}
		})
	}
}
