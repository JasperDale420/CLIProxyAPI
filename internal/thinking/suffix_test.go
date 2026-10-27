package thinking

import "testing"

func TestParseSuffix(t *testing.T) {
	tests := []struct {
		name, in, wantModel string
		wantHas             bool
		wantRaw             string
	}{
		{"no suffix", "gemini-2.5-pro", "gemini-2.5-pro", false, ""},
		{"numeric suffix", "claude-sonnet-4-5(16384)", "claude-sonnet-4-5", true, "16384"},
		{"level suffix", "gpt-5.2(high)", "gpt-5.2", true, "high"},
		{"unterminated paren is not a suffix", "model(high", "model(high", false, ""},
		{"text after paren is not a suffix", "model(high)x", "model(high)x", false, ""},
		{"last paren group wins", "a(b)(c)", "a(b)", true, "c"},
		{"empty suffix", "model()", "model", true, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSuffix(tt.in)
			if got.ModelName != tt.wantModel || got.HasSuffix != tt.wantHas || got.RawSuffix != tt.wantRaw {
				t.Fatalf("ParseSuffix(%q) = %+v", tt.in, got)
			}
		})
	}
}

func TestParseNumericSuffix(t *testing.T) {
	tests := []struct {
		in     string
		want   int
		wantOK bool
	}{
		{"8192", 8192, true},
		{"0", 0, true},
		{"08192", 8192, true},
		{"-1", 0, false},
		{"high", 0, false},
		{"", 0, false},
		{"99999999999999999999", 0, false},
	}
	for _, tt := range tests {
		got, ok := ParseNumericSuffix(tt.in)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("ParseNumericSuffix(%q) = (%d, %v), want (%d, %v)", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestParseSpecialSuffix(t *testing.T) {
	tests := []struct {
		in     string
		want   ThinkingMode
		wantOK bool
	}{
		{"none", ModeNone, true},
		{"NONE", ModeNone, true},
		{"auto", ModeAuto, true},
		{"-1", ModeAuto, true},
		{"high", ModeBudget, false},
		{"", ModeBudget, false},
	}
	for _, tt := range tests {
		got, ok := ParseSpecialSuffix(tt.in)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("ParseSpecialSuffix(%q) = (%v, %v), want (%v, %v)", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestParseLevelSuffix(t *testing.T) {
	tests := []struct {
		in     string
		want   ThinkingLevel
		wantOK bool
	}{
		{"minimal", LevelMinimal, true},
		{"low", LevelLow, true},
		{"Medium", LevelMedium, true},
		{"HIGH", LevelHigh, true},
		{"xhigh", LevelXHigh, true},
		{"max", LevelMax, true},
		{"none", "", false},
		{"auto", "", false},
		{"8192", "", false},
		{"ultra", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		got, ok := ParseLevelSuffix(tt.in)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("ParseLevelSuffix(%q) = (%q, %v), want (%q, %v)", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}
