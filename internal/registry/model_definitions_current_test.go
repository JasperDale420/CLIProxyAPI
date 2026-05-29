package registry

import "testing"

func TestCurrentCodexProModelsAreEmbedded(t *testing.T) {
	models := GetCodexProModels()
	for _, id := range []string{"gpt-5.4-mini", "gpt-5.5", "codex-auto-review"} {
		if !containsModelID(models, id) {
			t.Fatalf("expected Codex Pro model %q to be embedded", id)
		}
	}
}

func TestCurrentGeminiAndAntigravityModelsAreEmbedded(t *testing.T) {
	if !containsModelID(GetGeminiModels(), "gemini-3.5-flash") {
		t.Fatal("expected Gemini model gemini-3.5-flash to be embedded")
	}
	if !containsModelID(GetAntigravityModels(), "gemini-3.5-flash-low") {
		t.Fatal("expected Antigravity model gemini-3.5-flash-low to be embedded")
	}
}

func TestXAIStaticModelsAreAvailableByChannel(t *testing.T) {
	models := GetStaticModelDefinitionsByChannel("xai")
	if !containsModelID(models, "grok-4.3") {
		t.Fatal("expected xAI model grok-4.3 to be available")
	}
	if !containsModelID(models, "grok-imagine-image") {
		t.Fatal("expected xAI built-in model grok-imagine-image to be available")
	}
}

func containsModelID(models []*ModelInfo, id string) bool {
	for _, model := range models {
		if model != nil && model.ID == id {
			return true
		}
	}
	return false
}
