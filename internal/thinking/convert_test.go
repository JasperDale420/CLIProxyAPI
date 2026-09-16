package thinking

import (
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v7/internal/registry"
)

func TestConvertLevelToBudget(t *testing.T) {
	tests := []struct {
		name       string
		level      string
		wantBudget int
		wantOK     bool
	}{
		{name: "none", level: "none", wantBudget: 0, wantOK: true},
		{name: "auto", level: "auto", wantBudget: -1, wantOK: true},
		{name: "minimal", level: "minimal", wantBudget: 512, wantOK: true},
		{name: "low", level: "low", wantBudget: 1024, wantOK: true},
		{name: "medium", level: "medium", wantBudget: 8192, wantOK: true},
		{name: "high", level: "high", wantBudget: 24576, wantOK: true},
		{name: "xhigh", level: "xhigh", wantBudget: 32768, wantOK: true},
		{name: "max", level: "max", wantBudget: 128000, wantOK: true},
		{name: "case insensitive", level: "HIGH", wantBudget: 24576, wantOK: true},
		{name: "unknown level", level: "ultra", wantBudget: 0, wantOK: false},
		{name: "empty string", level: "", wantBudget: 0, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			budget, ok := ConvertLevelToBudget(tt.level)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if budget != tt.wantBudget {
				t.Errorf("budget = %d, want %d", budget, tt.wantBudget)
			}
		})
	}
}

func TestConvertBudgetToLevel(t *testing.T) {
	tests := []struct {
		name      string
		budget    int
		wantLevel string
		wantOK    bool
	}{
		{name: "invalid negative", budget: -2, wantLevel: "", wantOK: false},
		{name: "auto sentinel", budget: -1, wantLevel: string(LevelAuto), wantOK: true},
		{name: "zero is none", budget: 0, wantLevel: string(LevelNone), wantOK: true},
		{name: "one is minimal", budget: 1, wantLevel: string(LevelMinimal), wantOK: true},
		{name: "minimal upper bound", budget: 512, wantLevel: string(LevelMinimal), wantOK: true},
		{name: "just above minimal bound is low", budget: 513, wantLevel: string(LevelLow), wantOK: true},
		{name: "low upper bound", budget: 1024, wantLevel: string(LevelLow), wantOK: true},
		{name: "just above low bound is medium", budget: 1025, wantLevel: string(LevelMedium), wantOK: true},
		{name: "medium upper bound", budget: 8192, wantLevel: string(LevelMedium), wantOK: true},
		{name: "just above medium bound is high", budget: 8193, wantLevel: string(LevelHigh), wantOK: true},
		{name: "high upper bound", budget: 24576, wantLevel: string(LevelHigh), wantOK: true},
		{name: "just above high bound is xhigh", budget: 24577, wantLevel: string(LevelXHigh), wantOK: true},
		{name: "very large budget is xhigh", budget: 1000000, wantLevel: string(LevelXHigh), wantOK: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, ok := ConvertBudgetToLevel(tt.budget)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if level != tt.wantLevel {
				t.Errorf("level = %q, want %q", level, tt.wantLevel)
			}
		})
	}
}

func TestHasLevel(t *testing.T) {
	levels := []string{"low", " Medium ", "HIGH"}

	tests := []struct {
		name   string
		target string
		want   bool
	}{
		{name: "exact match", target: "low", want: true},
		{name: "case insensitive match", target: "medium", want: true},
		{name: "match ignores stored whitespace", target: "MEDIUM", want: true},
		{name: "match ignores case of stored value", target: "high", want: true},
		{name: "not present", target: "xhigh", want: false},
		{name: "empty target not present", target: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasLevel(levels, tt.target); got != tt.want {
				t.Errorf("HasLevel(%v, %q) = %v, want %v", levels, tt.target, got, tt.want)
			}
		})
	}

	if HasLevel(nil, "low") {
		t.Errorf("HasLevel(nil, ...) = true, want false")
	}
}

func TestMapToClaudeEffort(t *testing.T) {
	tests := []struct {
		name        string
		level       string
		supportsMax bool
		wantEffort  string
		wantOK      bool
	}{
		{name: "empty level", level: "", supportsMax: true, wantEffort: "", wantOK: false},
		{name: "minimal maps to low", level: "minimal", supportsMax: false, wantEffort: "low", wantOK: true},
		{name: "low passthrough", level: "low", supportsMax: false, wantEffort: "low", wantOK: true},
		{name: "medium passthrough", level: "medium", supportsMax: false, wantEffort: "medium", wantOK: true},
		{name: "high passthrough", level: "high", supportsMax: false, wantEffort: "high", wantOK: true},
		{name: "xhigh downgrades to high when max unsupported", level: "xhigh", supportsMax: false, wantEffort: "high", wantOK: true},
		{name: "xhigh maps to max when supported", level: "xhigh", supportsMax: true, wantEffort: "max", wantOK: true},
		{name: "max downgrades to high when max unsupported", level: "max", supportsMax: false, wantEffort: "high", wantOK: true},
		{name: "max maps to max when supported", level: "max", supportsMax: true, wantEffort: "max", wantOK: true},
		{name: "auto maps to high", level: "auto", supportsMax: true, wantEffort: "high", wantOK: true},
		{name: "unknown level", level: "ultra", supportsMax: true, wantEffort: "", wantOK: false},
		{name: "case insensitive and trimmed", level: " HIGH ", supportsMax: true, wantEffort: "high", wantOK: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			effort, ok := MapToClaudeEffort(tt.level, tt.supportsMax)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if effort != tt.wantEffort {
				t.Errorf("effort = %q, want %q", effort, tt.wantEffort)
			}
		})
	}
}

func TestDetectModelCapability(t *testing.T) {
	tests := []struct {
		name  string
		model *registry.ModelInfo
		want  ModelCapability
	}{
		{name: "nil model info is unknown", model: nil, want: CapabilityUnknown},
		{
			name:  "nil thinking support is none",
			model: &registry.ModelInfo{ID: "no-thinking", Thinking: nil},
			want:  CapabilityNone,
		},
		{
			name: "min/max without levels is budget only",
			model: &registry.ModelInfo{
				ID:       "budget-model",
				Thinking: &registry.ThinkingSupport{Min: 1024, Max: 128000},
			},
			want: CapabilityBudgetOnly,
		},
		{
			name: "levels without min/max is level only",
			model: &registry.ModelInfo{
				ID:       "level-model",
				Thinking: &registry.ThinkingSupport{Levels: []string{"low", "medium", "high"}},
			},
			want: CapabilityLevelOnly,
		},
		{
			name: "both min/max and levels is hybrid",
			model: &registry.ModelInfo{
				ID: "hybrid-model",
				Thinking: &registry.ThinkingSupport{
					Min:    128,
					Max:    20000,
					Levels: []string{"low", "medium", "high"},
				},
			},
			want: CapabilityHybrid,
		},
		{
			name: "thinking support present but empty is none",
			model: &registry.ModelInfo{
				ID:       "empty-support-model",
				Thinking: &registry.ThinkingSupport{},
			},
			want: CapabilityNone,
		},
		{
			name: "max alone without min still counts as budget",
			model: &registry.ModelInfo{
				ID:       "max-only-model",
				Thinking: &registry.ThinkingSupport{Max: 128000},
			},
			want: CapabilityBudgetOnly,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectModelCapability(tt.model); got != tt.want {
				t.Errorf("detectModelCapability() = %v, want %v", got, tt.want)
			}
		})
	}
}
