package thinking

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestGetThinkingText_ExtractsFromVariousFormats(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{
			name: "gemini-style direct text field",
			json: `{"thought":true,"text":"gemini thinking text"}`,
			want: "gemini thinking text",
		},
		{
			name: "simple thinking string",
			json: `{"thinking":"plain thinking text"}`,
			want: "plain thinking text",
		},
		{
			name: "wrapped thinking object with text field",
			json: `{"thinking":{"text":"wrapped text","cache_control":{"type":"ephemeral"}}}`,
			want: "wrapped text",
		},
		{
			name: "wrapped thinking object with inner thinking field",
			json: `{"thinking":{"thinking":"inner thinking text"}}`,
			want: "inner thinking text",
		},
		{
			name: "text field takes priority over thinking field",
			json: `{"text":"top-level text","thinking":"should be ignored"}`,
			want: "top-level text",
		},
		{
			name: "no thinking or text field returns empty",
			json: `{"other":"value"}`,
			want: "",
		},
		{
			name: "thinking object with neither text nor thinking inner field",
			json: `{"thinking":{"cache_control":{"type":"ephemeral"}}}`,
			want: "",
		},
		{
			name: "non-string text field is ignored",
			json: `{"text":123,"thinking":"fallback text"}`,
			want: "fallback text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part := gjson.Parse(tt.json)
			got := GetThinkingText(part)
			if got != tt.want {
				t.Errorf("GetThinkingText() = %q, want %q", got, tt.want)
			}
		})
	}
}
