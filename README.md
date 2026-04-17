# Empire AI Gateway

English | [中文](README_CN.md) | [日本語](README_JA.md)

## Quick Start

```bash
# Start the gateway
./cli-proxy-api

# Verify it's running
curl -H "Authorization: Bearer empire-ai-gateway-key" http://localhost:8002/v1/models
```

The gateway listens on **port 8002** and exposes an OpenAI-compatible API.

[![z.ai](https://assets.router-for.me/english-5-0.jpg)](https://z.ai/subscribe?ic=8JVLJQFSKB)

### API Key

GLM CODING PLAN is a subscription service designed for AI coding, starting at just $10/month. It provides access to their flagship GLM-4.7 & （GLM-5.1 Only Available  for Pro Users）model across 10+ popular AI coding tools (Claude Code, Cline, Roo Code, etc.), offering developers top-tier, fast, and stable coding experiences.

```
Authorization: Bearer empire-ai-gateway-key
```

### Provider Credentials

<table>
<tbody>
<tr>
<td width="180"><a href="https://www.packyapi.com/register?aff=cliproxyapi"><img src="./assets/packycode.png" alt="PackyCode" width="150"></a></td>
<td>Thanks to PackyCode for sponsoring this project! PackyCode is a reliable and efficient API relay service provider, offering relay services for Claude Code, Codex, Gemini, and more. PackyCode provides special discounts for our software users: register using <a href="https://www.packyapi.com/register?aff=cliproxyapi">this link</a> and enter the "cliproxyapi" promo code during recharge to get 10% off.</td>
</tr>
<tr>
<td width="180"><a href="https://www.aicodemirror.com/register?invitecode=TJNAIF"><img src="./assets/aicodemirror.png" alt="AICodeMirror" width="150"></a></td>
<td>Thanks to AICodeMirror for sponsoring this project! AICodeMirror provides official high-stability relay services for Claude Code / Codex / Gemini CLI, with enterprise-grade concurrency, fast invoicing, and 24/7 dedicated technical support. Claude Code / Codex / Gemini official channels at 38% / 2% / 9% of original price, with extra discounts on top-ups! AICodeMirror offers special benefits for CLIProxyAPI users: register via <a href="https://www.aicodemirror.com/register?invitecode=TJNAIF">this link</a> to enjoy 20% off your first top-up, and enterprise customers can get up to 25% off!</td>
</tr>
<tr>
<td width="180"><a href="https://shop.bmoplus.com/?utm_source=github"><img src="./assets/bmoplus.png" alt="BmoPlus" width="150"></a></td>
<td>Huge thanks to BmoPlus for sponsoring this project! BmoPlus is a highly reliable AI account provider built strictly for heavy AI users and developers. They offer rock-solid, ready-to-use accounts and official top-up services for ChatGPT Plus / ChatGPT Pro (Full Warranty) / Claude Pro / Super Grok / Gemini Pro. By registering and ordering through <a href="https://shop.bmoplus.com/?utm_source=github">BmoPlus - Premium AI Accounts & Top-ups</a>, users can unlock the mind-blowing rate of <b>10% of the official GPT subscription price (90% OFF)</b>!</td>
</tr>
<tr>
<td width="180"><a href="https://www.lingtrue.com/register"><img src="./assets/lingtrue.png" alt="LingtrueAPI" width="150"></a></td>
<td>Thanks to LingtrueAPI for its sponsorship of this project! LingtrueAPI is a global large - model API intermediary service platform that provides API calling services for various top - notch models such as Claude Code, Codex, and Gemini. It is committed to enabling users to connect to global AI capabilities at low cost and with high stability. LingtrueAPI offers special discounts to users of this software: register using <a href="https://www.lingtrue.com/register">this link</a>, and enter the promo code "LingtrueAPI" when making the first recharge to enjoy a 10% discount.</td>
</tr>
</tbody>
</table>

```
http://localhost:8002/management.html
```

**Management Key:** `empire`

From the management UI you can:

- Add **Antigravity** accounts (Google OAuth → Gemini, Claude, GPT models)
- Add **Codex** accounts (Apple OAuth → OpenAI GPT models)
- Add **Claude Code** accounts (Anthropic OAuth)
- Add **Gemini CLI** accounts (Google OAuth)
- View usage statistics and account status

To add accounts via CLI instead:

```bash
./cli-proxy-api --antigravity-login    # Antigravity (Gemini/Claude/GPT)
./cli-proxy-api --codex-login          # Codex (OpenAI GPT)
./cli-proxy-api --claude-login         # Claude Code
./cli-proxy-api --gemini-login         # Gemini CLI
```

## Making Requests

### Endpoint

```
POST http://localhost:8002/v1/chat/completions
```

When you need the request/response shape of a specific backend family, use the provider-specific paths instead of the merged `/v1/...` endpoints:

- Use `/api/provider/{provider}/v1/messages` for messages-style backends.
- Use `/api/provider/{provider}/v1beta/models/...` for model-scoped generate endpoints.
- Use `/api/provider/{provider}/v1/chat/completions` for chat-completions backends.

These routes help you select the protocol surface, but they do not by themselves guarantee a unique inference executor when the same client-visible model name is reused across multiple backends. Inference routing is still resolved from the request model/alias. For strict backend pinning, use unique aliases, prefixes, or otherwise avoid overlapping client-visible model names.

**→ [Complete Amp CLI Integration Guide](https://help.router-for.me/agent-client/amp-cli.html)**

```bash
curl -X POST http://localhost:8002/v1/chat/completions \
  -H "Authorization: Bearer empire-ai-gateway-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gemini-3-pro-high",
    "messages": [{"role": "user", "content": "Hello"}]
  }'
```

### Streaming

```bash
curl -N -X POST http://localhost:8002/v1/chat/completions \
  -H "Authorization: Bearer empire-ai-gateway-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gemini-3-pro-high",
    "messages": [{"role": "user", "content": "Hello"}],
    "stream": true
  }'
```

### Python

```python
import openai

client = openai.OpenAI(
    base_url="http://localhost:8002/v1",
    api_key="empire-ai-gateway-key",
)

response = client.chat.completions.create(
    model="gemini-3-pro-high",
    messages=[{"role": "user", "content": "Hello"}],
)
print(response.choices[0].message.content)
```

### Any OpenAI-Compatible Client

Point any OpenAI-compatible SDK or tool at the gateway:

| Setting    | Value                              |
|------------|------------------------------------|
| Base URL   | `http://localhost:8002/v1`         |
| API Key    | `empire-ai-gateway-key`           |
| Model      | Any model from the table below     |

## Available Models

List models dynamically:

```bash
curl -H "Authorization: Bearer empire-ai-gateway-key" http://localhost:8002/v1/models
```

### Antigravity Provider

Authenticated via Google OAuth. Routes to Google, Anthropic, and OpenAI models through the Antigravity platform.

| Model ID | Description |
|----------|-------------|
| `gemini-3-pro-high` | Gemini 3 Pro (high quality) |
| `gemini-3-pro-image` | Gemini 3 Pro with image generation |
| `gemini-3-flash` | Gemini 3 Flash (fast) |
| `gemini-2.5-flash` | Gemini 2.5 Flash |
| `gemini-2.5-flash-lite` | Gemini 2.5 Flash Lite (fastest) |
| `claude-sonnet-4-6` | Claude Sonnet 4.6 |
| `claude-sonnet-4-5` | Claude Sonnet 4.5 |
| `claude-sonnet-4-5-thinking` | Claude Sonnet 4.5 with extended thinking |
| `claude-opus-4-6-thinking` | Claude Opus 4.6 with extended thinking |
| `gpt-oss-120b-medium` | GPT OSS 120B |

### Codex Provider (OpenAI)

Authenticated via Apple OAuth. Routes to OpenAI GPT models.

| Model ID | Description |
|----------|-------------|
| `gpt-5` | GPT-5 |
| `gpt-5-codex` | GPT-5 Codex |
| `gpt-5-codex-mini` | GPT-5 Codex Mini |
| `gpt-5.1` | GPT-5.1 |
| `gpt-5.1-codex` | GPT-5.1 Codex |
| `gpt-5.1-codex-mini` | GPT-5.1 Codex Mini |
| `gpt-5.1-codex-max` | GPT-5.1 Codex Max |
| `gpt-5.2` | GPT-5.2 |
| `gpt-5.2-codex` | GPT-5.2 Codex |
| `gpt-5.3-codex` | GPT-5.3 Codex |
| `gpt-5.3-codex-spark` | GPT-5.3 Codex Spark |

### Zhipu GLM Coding Plan

Routes through the official GLM coding endpoint. This is the correct path for Z.ai / BigModel coding-plan API keys.

| Model ID | Description |
|----------|-------------|
| `glm-4.5` | GLM 4.5 |
| `glm-4.5-air` | GLM 4.5 Air |
| `glm-4.6` | GLM 4.6 |
| `glm-4.7` | GLM 4.7 |
| `glm-5` | GLM 5 |

Flash and vision variants are not advertised in this gateway config because the coding endpoint does not expose them in `/models`.

### MackingJAI (ChatGPT Desktop App)

Routes through the ChatGPT macOS desktop app via [MackingJAI](https://github.com/0ssamaak0/MackingJAI). **No rate limits** — uses your ChatGPT subscription directly.

> **Requires:** MackingJAI.app running + [Apple Shortcut installed](https://www.icloud.com/shortcuts/753cd6efc8fb49918817e107f12a0420) + ChatGPT desktop app running.

| Model ID | Description |
|----------|-------------|
| `chatgpt` | GPT-5 via ChatGPT desktop (no rate limits) |
| `chatgpt-unlimited` | GPT-5 via ChatGPT desktop (alias) |

**How model routing works:** ChatGPT's desktop app auto-routes all requests to GPT-5. There is no way to select specific model variants (e.g. `fast`, `thinking`) or legacy models (o3, gpt-4o) — the ChatGPT Shortcut integration does not expose model selection. Both `chatgpt` and `chatgpt-unlimited` are aliases that map to `GPT-5` and behave identically.

#### MackingJAI vs Codex OAuth

| | MackingJAI (`chatgpt`) | Codex OAuth (`gpt-5`, `gpt-5-codex`, etc.) |
|---|---|---|
| **Rate limits** | None (uses subscription) | Standard API limits |
| **Model control** | Auto-routed to GPT-5 only | Exact variant selection |
| **Parameters** | No temp/top_p/penalties | Full parameter control |
| **Streaming** | Not supported | Supported |
| **Image inputs** | Not supported | Supported |
| **Function calling** | Not supported | Supported |
| **Speed** | Slower (Shortcut overhead) | Faster (direct API) |
| **Best for** | Bulk/batch work, avoiding rate limits | Production use, precise control |

**When to use `chatgpt`:** Batch processing, bulk requests, or any workload where API rate limits are the bottleneck and you don't need streaming, function calling, or fine-grained model control.

**When to use Codex models:** Production services, real-time applications, or any workload requiring streaming, exact model selection, tool use, or parameter tuning.

> **Note:** Model availability depends on which provider accounts are authenticated. Run `/v1/models` to see your current list.

### [All API Hub](https://github.com/qixing-jk/all-api-hub)

Browser extension for one-stop management of New API-compatible relay site accounts, featuring balance and usage dashboards, auto check-in, one-click key export to common apps, in-page API availability testing, and channel/model sync and redirection. It integrates with CLIProxyAPI through the Management API for one-click provider import and config sync.

### [Shadow AI](https://github.com/HEUDavid/shadow-ai)

Shadow AI is an AI assistant tool designed specifically for restricted environments. It provides a stealthy operation
mode without windows or traces, and enables cross-device AI Q&A interaction and control via the local area network (
LAN). Essentially, it is an automated collaboration layer of "screen/audio capture + AI inference + low-friction delivery",
helping users to immersively use AI assistants across applications on controlled devices or in restricted environments.

> [!NOTE]
> If you developed a project based on CLIProxyAPI, please open a PR to add it to this list.

- **OpenAI-compatible API** — works with any OpenAI SDK or client
- **Multi-provider routing** — single endpoint for Gemini, Claude, and GPT models
- **Streaming & non-streaming** — both response modes supported
- **Function calling / tools** — full tool use support
- **Multimodal** — text and image inputs
- **Multi-account load balancing** — round-robin across credentials
- **Auto-retry** — 3 retries on 403/408/500/502/503/504
- **Quota management** — auto-switches accounts on 429, falls back to preview models
- **Hot-reload config** — edit `config.yaml` without restarting

## Configuration

Configuration lives in `config.yaml` (hot-reloaded on save). Key settings:

| Setting | Default | Description |
|---------|---------|-------------|
| `port` | `8002` | HTTP listen port |
| `api-keys` | `["empire-ai-gateway-key"]` | Keys that clients must provide |
| `auth-dir` | `~/.cli-proxy-api` | Where OAuth credentials are stored |
| `request-retry` | `3` | Retry count on server errors |
| `quota-exceeded.switch-project` | `true` | Auto-switch accounts on 429 |
| `quota-exceeded.switch-preview-model` | `true` | Fall back to preview models |
| `routing.strategy` | `round-robin` | Credential selection strategy |
| `debug` | `false` | Verbose debug logging |

### [OmniRoute](https://github.com/diegosouzapw/OmniRoute)

Never stop coding. Smart routing to FREE & low-cost AI models with automatic fallback.

OmniRoute is an AI gateway for multi-provider LLMs: an OpenAI-compatible endpoint with smart routing, load balancing, retries, and fallbacks. Add policies, rate limits, caching, and observability for reliable, cost-aware inference.

> [!NOTE]
> If you have developed a port of CLIProxyAPI or a project inspired by it, please open a PR to add it to this list.

## Management

### Web UI

```
http://localhost:8002/management.html
```

Login with management key `empire` to:

- Add/remove provider accounts
- View usage statistics
- Monitor account health

### Management API

```bash
# Get config
curl -H "Authorization: Bearer empire" http://localhost:8002/v0/management/config

# List auth files
curl -H "Authorization: Bearer empire" http://localhost:8002/v0/management/auth-files

# View usage
curl -H "Authorization: Bearer empire" http://localhost:8002/v0/management/usage
```

## Troubleshooting

### `unknown provider for model X`

You're using a model name that doesn't match any configured provider. Check available models with `/v1/models`.

### `502 Bad Gateway`

The gateway can't route to a provider. Verify credentials are valid in the management UI.

### No response / hanging requests

Use non-streaming mode or Python `urllib`/`requests` for testing. Some `curl` configurations wait for streaming chunks indefinitely. Always add `--max-time` for safety.

### Refresh credentials

Credentials auto-refresh every 15 minutes. To force a refresh, restart the gateway or re-authenticate via the management UI.

## Upstream

Based on [CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI) v6.8.51 by [router-for-me](https://github.com/router-for-me).
Full documentation: [https://help.router-for.me/](https://help.router-for.me/)
