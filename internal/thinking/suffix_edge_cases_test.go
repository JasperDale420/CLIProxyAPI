package thinking

import "testing"

func TestParseNumericSuffix_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		wantBudget int
		wantOK     bool
	}{
		{"empty string is invalid", "", 0, false},
		{"positive integer", "8192", 8192, true},
		{"zero is valid", "0", 0, true},
		{"leading zeros accepted", "08192", 8192, true},
		{"negative number is invalid", "-1", 0, false},
		{"non-numeric string is invalid", "high", 0, false},
		{"overflow on 64-bit is invalid", "9223372036854775808", 0, false},
		{"whitespace is invalid", " 100", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			budget, ok := ParseNumericSuffix(tt.raw)
			if ok != tt.wantOK {
				t.Fatalf("ParseNumericSuffix(%q) ok = %v, want %v", tt.raw, ok, tt.wantOK)
			}
			if ok && budget != tt.wantBudget {
				t.Fatalf("ParseNumericSuffix(%q) budget = %d, want %d", tt.raw, budget, tt.wantBudget)
			}
		})
	}
}

func TestParseSpecialSuffix_CaseInsensitiveMatching(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		wantMode ThinkingMode
		wantOK   bool
	}{
		{"none lowercase", "none", ModeNone, true},
		{"NONE uppercase", "NONE", ModeNone, true},
		{"auto lowercase", "auto", ModeAuto, true},
		{"Auto mixed case", "Auto", ModeAuto, true},
		{"-1 maps to auto", "-1", ModeAuto, true},
		{"empty string is invalid", "", ModeBudget, false},
		{"unrelated numeric value is invalid", "8192", ModeBudget, false},
		{"unknown word is invalid", "sometimes", ModeBudget, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, ok := ParseSpecialSuffix(tt.raw)
			if ok != tt.wantOK {
				t.Fatalf("ParseSpecialSuffix(%q) ok = %v, want %v", tt.raw, ok, tt.wantOK)
			}
			if mode != tt.wantMode {
				t.Fatalf("ParseSpecialSuffix(%q) mode = %v, want %v", tt.raw, mode, tt.wantMode)
			}
		})
	}
}

func TestParseLevelSuffix_RejectsSpecialAndNumericValues(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		wantLevel ThinkingLevel
		wantOK    bool
	}{
		{"high level matches", "high", LevelHigh, true},
		{"HIGH uppercase matches", "HIGH", LevelHigh, true},
		{"xhigh level matches", "xhigh", LevelXHigh, true},
		{"none is a special value, not a level", "none", "", false},
		{"auto is a special value, not a level", "auto", "", false},
		{"numeric suffix is not a level", "8192", "", false},
		{"unknown level string", "ultra", "", false},
		{"empty string is invalid", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, ok := ParseLevelSuffix(tt.raw)
			if ok != tt.wantOK {
				t.Fatalf("ParseLevelSuffix(%q) ok = %v, want %v", tt.raw, ok, tt.wantOK)
			}
			if level != tt.wantLevel {
				t.Fatalf("ParseLevelSuffix(%q) level = %v, want %v", tt.raw, level, tt.wantLevel)
			}
		})
	}
}
