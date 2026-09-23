package diff

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/router-for-me/CLIProxyAPI/v7/internal/config"
)

// ComputeOpenAICompatModelsHash returns a stable hash for OpenAI-compat models.
// Used to detect model list changes during hot reload.
func ComputeOpenAICompatModelsHash(models []config.OpenAICompatibilityModel) string {
	keys := normalizeModelPairs(func(out func(key string)) {
		for _, model := range models {
			name := strings.TrimSpace(model.Name)
			alias := strings.TrimSpace(model.Alias)
			if name == "" && alias == "" {
				continue
			}
			out(strings.ToLower(name) + "|" + strings.ToLower(alias) + "|" + fmt.Sprintf("image=%t", model.Image))
		}
	})
	return hashJoined(keys)
}

// namedAliasModel is implemented by the config model types that carry a plain
// upstream Name and a client-facing Alias, letting their hash computation
// share one implementation instead of four copy-pasted loops.
type namedAliasModel interface {
	GetName() string
	GetAlias() string
}

// computeNamedAliasModelsHash returns a stable hash over the normalized
// name/alias pairs of models. Shared by the Compute*ModelsHash helpers below
// whose models differ only in type, not in which fields feed the hash.
func computeNamedAliasModelsHash[T namedAliasModel](models []T) string {
	keys := normalizeModelPairs(func(out func(key string)) {
		for _, model := range models {
			name := strings.TrimSpace(model.GetName())
			alias := strings.TrimSpace(model.GetAlias())
			if name == "" && alias == "" {
				continue
			}
			out(strings.ToLower(name) + "|" + strings.ToLower(alias))
		}
	})
	return hashJoined(keys)
}

// ComputeVertexCompatModelsHash returns a stable hash for Vertex-compatible models.
func ComputeVertexCompatModelsHash(models []config.VertexCompatModel) string {
	return computeNamedAliasModelsHash(models)
}

// ComputeClaudeModelsHash returns a stable hash for Claude model aliases.
func ComputeClaudeModelsHash(models []config.ClaudeModel) string {
	return computeNamedAliasModelsHash(models)
}

// ComputeCodexModelsHash returns a stable hash for Codex model aliases.
func ComputeCodexModelsHash(models []config.CodexModel) string {
	return computeNamedAliasModelsHash(models)
}

// ComputeGeminiModelsHash returns a stable hash for Gemini model aliases.
func ComputeGeminiModelsHash(models []config.GeminiModel) string {
	return computeNamedAliasModelsHash(models)
}

// ComputeExcludedModelsHash returns a normalized hash for excluded model lists.
func ComputeExcludedModelsHash(excluded []string) string {
	if len(excluded) == 0 {
		return ""
	}
	normalized := make([]string, 0, len(excluded))
	for _, entry := range excluded {
		if trimmed := strings.TrimSpace(entry); trimmed != "" {
			normalized = append(normalized, strings.ToLower(trimmed))
		}
	}
	if len(normalized) == 0 {
		return ""
	}
	sort.Strings(normalized)
	data, _ := json.Marshal(normalized)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func normalizeModelPairs(collect func(out func(key string))) []string {
	seen := make(map[string]struct{})
	keys := make([]string, 0)
	collect(func(key string) {
		if _, exists := seen[key]; exists {
			return
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	})
	if len(keys) == 0 {
		return nil
	}
	sort.Strings(keys)
	return keys
}

func hashJoined(keys []string) string {
	if len(keys) == 0 {
		return ""
	}
	sum := sha256.Sum256([]byte(strings.Join(keys, "\n")))
	return hex.EncodeToString(sum[:])
}
