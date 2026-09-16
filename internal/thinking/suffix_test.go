package thinking

import "testing"

func TestParseSuffix(t *testing.T) {
	tests := []struct {
		name          string
		model         string
		wantModelName string
		wantHasSuffix bool
		wantRawSuffix string
	}{
		{
			name:          "numeric suffix",
			model:         "claude-sonnet-4-5(16384)",
			wantModelName: "claude-sonnet-4-5",
			wantHasSuffix: true,
			wantRawSuffix: "16384",
		},
		{
			name:          "level suffix",
			model:         "gpt-5.2(high)",
			wantModelName: "gpt-5.2",
			wantHasSuffix: true,
			wantRawSuffix: "high",
		},
		{
			name:          "no suffix",
			model:         "gemini-2.5-pro",
			wantModelName: "gemini-2.5-pro",
			wantHasSuffix: false,
			wantRawSuffix: "",
		},
		{
			name:          "empty parens",
			model:         "gpt-5.2()",
			wantModelName: "gpt-5.2",
			wantHasSuffix: true,
			wantRawSuffix: "",
		},
		{
			name:          "unclosed paren is not a suffix",
			model:         "gpt-5.2(high",
			wantModelName: "gpt-5.2(high",
			wantHasSuffix: false,
			wantRawSuffix: "",
		},
		{
			name: "opening paren with no closing paren anywhere but string still ends " +
				"with a stray char is not a suffix",
			model:         "gpt-5.2(high)x",
			wantModelName: "gpt-5.2(high)x",
			wantHasSuffix: false,
			wantRawSuffix: "",
		},
		{
			name:          "nested parens use last opening paren",
			model:         "model(outer(inner))",
			wantModelName: "model(outer",
			wantHasSuffix: true,
			wantRawSuffix: "inner)",
		},
		{
			name:          "model name itself contains a matched paren pair before the suffix",
			model:         "model(v1)(high)",
			wantModelName: "model(v1)",
			wantHasSuffix: true,
			wantRawSuffix: "high",
		},
		{
			name:          "empty string",
			model:         "",
			wantModelName: "",
			wantHasSuffix: false,
			wantRawSuffix: "",
		},
		{
			name:          "only parens",
			model:         "()",
			wantModelName: "",
			wantHasSuffix: true,
			wantRawSuffix: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSuffix(tt.model)
			if got.ModelName != tt.wantModelName {
				t.Errorf("ModelName = %q, want %q", got.ModelName, tt.wantModelName)
			}
			if got.HasSuffix != tt.wantHasSuffix {
				t.Errorf("HasSuffix = %v, want %v", got.HasSuffix, tt.wantHasSuffix)
			}
			if got.RawSuffix != tt.wantRawSuffix {
				t.Errorf("RawSuffix = %q, want %q", got.RawSuffix, tt.wantRawSuffix)
			}
		})
	}
}

func TestParseNumericSuffix(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		wantBudget int
		wantOK     bool
	}{
		{name: "simple positive", raw: "8192", wantBudget: 8192, wantOK: true},
		{name: "zero", raw: "0", wantBudget: 0, wantOK: true},
		{name: "leading zeros accepted", raw: "08192", wantBudget: 8192, wantOK: true},
		{name: "negative is rejected", raw: "-1", wantBudget: 0, wantOK: false},
		{name: "non-numeric", raw: "high", wantBudget: 0, wantOK: false},
		{name: "empty string", raw: "", wantBudget: 0, wantOK: false},
		{name: "overflow on 64-bit int", raw: "9223372036854775808", wantBudget: 0, wantOK: false},
		{name: "whitespace is not numeric", raw: " 8192", wantBudget: 0, wantOK: false},
		{name: "float is not numeric", raw: "8192.0", wantBudget: 0, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			budget, ok := ParseNumericSuffix(tt.raw)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if budget != tt.wantBudget {
				t.Errorf("budget = %d, want %d", budget, tt.wantBudget)
			}
		})
	}
}

func TestParseSpecialSuffix(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		wantMode ThinkingMode
		wantOK   bool
	}{
		{name: "none lowercase", raw: "none", wantMode: ModeNone, wantOK: true},
		{name: "none uppercase is case-insensitive", raw: "NONE", wantMode: ModeNone, wantOK: true},
		{name: "auto", raw: "auto", wantMode: ModeAuto, wantOK: true},
		{name: "auto mixed case", raw: "Auto", wantMode: ModeAuto, wantOK: true},
		{name: "numeric -1 means auto", raw: "-1", wantMode: ModeAuto, wantOK: true},
		{name: "unrelated numeric value", raw: "8192", wantMode: ModeBudget, wantOK: false},
		{name: "unrelated level value", raw: "high", wantMode: ModeBudget, wantOK: false},
		{name: "empty string", raw: "", wantMode: ModeBudget, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, ok := ParseSpecialSuffix(tt.raw)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if mode != tt.wantMode {
				t.Errorf("mode = %v, want %v", mode, tt.wantMode)
			}
		})
	}
}

func TestParseLevelSuffix(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		wantLevel ThinkingLevel
		wantOK    bool
	}{
		{name: "minimal", raw: "minimal", wantLevel: LevelMinimal, wantOK: true},
		{name: "low", raw: "low", wantLevel: LevelLow, wantOK: true},
		{name: "medium", raw: "medium", wantLevel: LevelMedium, wantOK: true},
		{name: "high", raw: "high", wantLevel: LevelHigh, wantOK: true},
		{name: "xhigh", raw: "xhigh", wantLevel: LevelXHigh, wantOK: true},
		{name: "max", raw: "max", wantLevel: LevelMax, wantOK: true},
		{name: "case insensitive", raw: "HIGH", wantLevel: LevelHigh, wantOK: true},
		{name: "special value none is rejected", raw: "none", wantLevel: "", wantOK: false},
		{name: "special value auto is rejected", raw: "auto", wantLevel: "", wantOK: false},
		{name: "numeric value is rejected", raw: "8192", wantLevel: "", wantOK: false},
		{name: "unknown level is rejected", raw: "ultra", wantLevel: "", wantOK: false},
		{name: "empty string", raw: "", wantLevel: "", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, ok := ParseLevelSuffix(tt.raw)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if level != tt.wantLevel {
				t.Errorf("level = %q, want %q", level, tt.wantLevel)
			}
		})
	}
}
