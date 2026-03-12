// Package openai implements thinking configuration for OpenAI/Codex models.
//
// OpenAI models use the reasoning_effort format with discrete levels
// (low/medium/high). Some models support xhigh and none levels.
// See: _bmad-output/planning-artifacts/architecture.md#Epic-8
package openai

import (
	"strings"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/registry"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/thinking"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Applier implements thinking.ProviderApplier for OpenAI models.
//
// OpenAI-specific behavior:
//   - Output format: reasoning_effort (string: low/medium/high/xhigh)
//   - Level-only mode: no numeric budget support
//   - Some models support ZeroAllowed (gpt-5.1, gpt-5.2)
type Applier struct{}

var _ thinking.ProviderApplier = (*Applier)(nil)

// NewApplier creates a new OpenAI thinking applier.
func NewApplier() *Applier {
	return &Applier{}
}

func init() {
	thinking.RegisterProvider("openai", NewApplier())
}

// Apply applies thinking configuration to OpenAI request body.
//
// Expected output format:
//
//	{
//	  "reasoning_effort": "high"
//	}
func (a *Applier) Apply(body []byte, config thinking.ThinkingConfig, modelInfo *registry.ModelInfo) ([]byte, error) {
	if thinking.IsUserDefinedModel(modelInfo) {
		return applyCompatibleOpenAI(body, config)
	}
	if modelInfo.Thinking == nil {
		return body, nil
	}

	// Only handle ModeLevel and ModeNone; other modes pass through unchanged.
	if config.Mode != thinking.ModeLevel && config.Mode != thinking.ModeNone {
		return body, nil
	}

	if len(body) == 0 || !gjson.ValidBytes(body) {
		body = []byte(`{}`)
	}

	if config.Mode == thinking.ModeLevel {
		result, _ := sjson.SetBytes(body, "reasoning_effort", string(config.Level))
		return result, nil
	}

	effort := ""
	support := modelInfo.Thinking
	if config.Budget == 0 {
		if support.ZeroAllowed || thinking.HasLevel(support.Levels, string(thinking.LevelNone)) {
			effort = string(thinking.LevelNone)
		}
	}
	if effort == "" && config.Level != "" {
		effort = string(config.Level)
	}
	if effort == "" && len(support.Levels) > 0 {
		// OpenAI reasoning_effort natively supports low/medium/high. When the target
		// model exposes "minimal" as its first level but ZeroAllowed=false, skip
		// "minimal" (a sub-low tier) and use the first level at or above "low" so
		// the outgoing value is always a recognised OpenAI effort string.
		effort = support.Levels[0]
		if strings.EqualFold(effort, string(thinking.LevelMinimal)) && !support.ZeroAllowed {
			for _, l := range support.Levels[1:] {
				if !strings.EqualFold(l, string(thinking.LevelMinimal)) {
					effort = l
					break
				}
			}
		}
	}
	if effort == "" {
		return body, nil
	}

	result, _ := sjson.SetBytes(body, "reasoning_effort", effort)
	return result, nil
}

func applyCompatibleOpenAI(body []byte, config thinking.ThinkingConfig) ([]byte, error) {
	if len(body) == 0 || !gjson.ValidBytes(body) {
		body = []byte(`{}`)
	}

	var effort string
	switch config.Mode {
	case thinking.ModeLevel:
		if config.Level == "" {
			return body, nil
		}
		// OpenAI reasoning_effort supports low/medium/high only. Clamp supra-high levels
		// (xhigh, max) down to high, and map the special "auto" level to medium.
		level := string(config.Level)
		switch {
		case strings.EqualFold(level, string(thinking.LevelXHigh)), strings.EqualFold(level, string(thinking.LevelMax)):
			effort = string(thinking.LevelHigh)
		case strings.EqualFold(level, string(thinking.LevelAuto)):
			effort = string(thinking.LevelMedium)
		default:
			effort = level
		}
	case thinking.ModeNone:
		effort = string(thinking.LevelNone)
		if config.Level != "" {
			effort = string(config.Level)
		}
	case thinking.ModeAuto:
		// OpenAI reasoning_effort does not support "auto"; use medium as the default
		// mid-range level for user-defined models.
		effort = string(thinking.LevelMedium)
	case thinking.ModeBudget:
		// Budget mode: convert budget to level using threshold mapping.
		// OpenAI only supports low/medium/high; clamp anything above "high"
		// (e.g. xhigh, max) down to "high".
		level, ok := thinking.ConvertBudgetToLevel(config.Budget)
		if !ok {
			return body, nil
		}
		if strings.EqualFold(level, string(thinking.LevelXHigh)) || strings.EqualFold(level, string(thinking.LevelMax)) {
			level = string(thinking.LevelHigh)
		}
		effort = level
	default:
		return body, nil
	}

	result, _ := sjson.SetBytes(body, "reasoning_effort", effort)
	return result, nil
}
