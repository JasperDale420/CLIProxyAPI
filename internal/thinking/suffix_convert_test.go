package thinking

import "testing"

func TestParseSuffix(t *testing.T) {
	tests := []struct {
		name      string
		model     string
		wantModel string
		wantHas   bool
		wantRaw   string
	}{
		{"no suffix", "gemini-2.5-pro", "gemini-2.5-pro", false, ""},
		{"numeric suffix", "claude-sonnet-4-5(16384)", "claude-sonnet-4-5", true, "16384"},
		{"level suffix", "gpt-5.2(high)", "gpt-5.2", true, "high"},
		{"unclosed parenthesis is not a suffix", "gpt-5.2(high", "gpt-5.2(high", false, ""},
		{"closing parenthesis without opener is not a suffix", "gpt-5.2)", "gpt-5.2)", false, ""},
		{"last parenthesis group wins", "model(a)(b)", "model(a)", true, "b"},
		{"empty suffix still counts as a suffix", "model()", "model", true, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSuffix(tt.model)
			if got.ModelName != tt.wantModel || got.HasSuffix != tt.wantHas || got.RawSuffix != tt.wantRaw {
				t.Fatalf("ParseSuffix(%q) = %+v, want model=%q has=%v raw=%q", tt.model, got, tt.wantModel, tt.wantHas, tt.wantRaw)
			}
		})
	}
}

func TestParseNumericSuffix(t *testing.T) {
	tests := []struct {
		raw    string
		want   int
		wantOK bool
	}{
		{"8192", 8192, true},
		{"0", 0, true},
		{"08192", 8192, true},
		{"-1", 0, false},
		{"-5", 0, false},
		{"high", 0, false},
		{"", 0, false},
		{"12.5", 0, false},
		{"99999999999999999999", 0, false},
	}
	for _, tt := range tests {
		got, ok := ParseNumericSuffix(tt.raw)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("ParseNumericSuffix(%q) = (%d, %v), want (%d, %v)", tt.raw, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestParseSpecialSuffix(t *testing.T) {
	tests := []struct {
		raw      string
		wantMode ThinkingMode
		wantOK   bool
	}{
		{"none", ModeNone, true},
		{"NONE", ModeNone, true},
		{"auto", ModeAuto, true},
		{"Auto", ModeAuto, true},
		{"-1", ModeAuto, true},
		{"high", ModeBudget, false},
		{"0", ModeBudget, false},
		{"", ModeBudget, false},
	}
	for _, tt := range tests {
		mode, ok := ParseSpecialSuffix(tt.raw)
		if mode != tt.wantMode || ok != tt.wantOK {
			t.Errorf("ParseSpecialSuffix(%q) = (%v, %v), want (%v, %v)", tt.raw, mode, ok, tt.wantMode, tt.wantOK)
		}
	}
}

func TestParseLevelSuffix(t *testing.T) {
	valid := map[string]ThinkingLevel{
		"minimal": LevelMinimal,
		"low":     LevelLow,
		"medium":  LevelMedium,
		"high":    LevelHigh,
		"HIGH":    LevelHigh,
		"xhigh":   LevelXHigh,
		"max":     LevelMax,
	}
	for raw, want := range valid {
		got, ok := ParseLevelSuffix(raw)
		if !ok || got != want {
			t.Errorf("ParseLevelSuffix(%q) = (%q, %v), want (%q, true)", raw, got, ok, want)
		}
	}
	// Special and numeric values belong to the other parsers.
	for _, raw := range []string{"", "none", "auto", "8192", "ultra"} {
		if got, ok := ParseLevelSuffix(raw); ok || got != "" {
			t.Errorf("ParseLevelSuffix(%q) = (%q, %v), want (\"\", false)", raw, got, ok)
		}
	}
}

func TestConvertLevelToBudget(t *testing.T) {
	tests := []struct {
		level  string
		want   int
		wantOK bool
	}{
		{"none", 0, true},
		{"auto", -1, true},
		{"minimal", 512, true},
		{"low", 1024, true},
		{"medium", 8192, true},
		{"HIGH", 24576, true},
		{"xhigh", 32768, true},
		{"max", 128000, true},
		{"ultra", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		got, ok := ConvertLevelToBudget(tt.level)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("ConvertLevelToBudget(%q) = (%d, %v), want (%d, %v)", tt.level, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestConvertBudgetToLevelThresholdBoundaries(t *testing.T) {
	tests := []struct {
		budget int
		want   string
		wantOK bool
	}{
		{-2, "", false},
		{-1, "auto", true},
		{0, "none", true},
		{1, "minimal", true},
		{512, "minimal", true},
		{513, "low", true},
		{1024, "low", true},
		{1025, "medium", true},
		{8192, "medium", true},
		{8193, "high", true},
		{24576, "high", true},
		{24577, "xhigh", true},
		{1000000, "xhigh", true},
	}
	for _, tt := range tests {
		got, ok := ConvertBudgetToLevel(tt.budget)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("ConvertBudgetToLevel(%d) = (%q, %v), want (%q, %v)", tt.budget, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestHasLevelIgnoresCaseAndWhitespace(t *testing.T) {
	levels := []string{"low", " High ", "MEDIUM"}
	for _, target := range []string{"low", "high", "Medium"} {
		if !HasLevel(levels, target) {
			t.Errorf("HasLevel(%v, %q) = false, want true", levels, target)
		}
	}
	if HasLevel(levels, "max") {
		t.Error("HasLevel reported a level that is not present")
	}
	if HasLevel(nil, "low") {
		t.Error("HasLevel on an empty slice must be false")
	}
}

func TestMapToClaudeEffort(t *testing.T) {
	tests := []struct {
		level       string
		supportsMax bool
		want        string
		wantOK      bool
	}{
		{"minimal", false, "low", true},
		{"low", false, "low", true},
		{" Medium ", false, "medium", true},
		{"high", true, "high", true},
		{"xhigh", true, "max", true},
		{"max", true, "max", true},
		{"xhigh", false, "high", true},
		{"max", false, "high", true},
		{"auto", false, "high", true},
		{"", false, "", false},
		{"none", false, "", false},
		{"ultra", true, "", false},
	}
	for _, tt := range tests {
		got, ok := MapToClaudeEffort(tt.level, tt.supportsMax)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("MapToClaudeEffort(%q, %v) = (%q, %v), want (%q, %v)", tt.level, tt.supportsMax, got, ok, tt.want, tt.wantOK)
		}
	}
}
