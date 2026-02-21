# Empire AI Gateway

Unified LLM proxy powered by [CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI) (v6.8.21). Provides a single OpenAI-compatible API endpoint for routing requests to multiple LLM providers using OAuth-authenticated CLI credentials.

## Quick Start

```bash
# Start the gateway
./cli-proxy-api

# Verify it's running
curl -H "Authorization: Bearer empire-ai-gateway-key" http://localhost:8002/v1/models
```

The gateway listens on **port 8002** and exposes an OpenAI-compatible API.

## Authentication

### API Key

All requests require the API key in the `Authorization` header:

```
Authorization: Bearer empire-ai-gateway-key
```

### Provider Credentials

The gateway proxies to LLM providers using OAuth credentials stored in `~/.cli-proxy-api/`. To add or manage accounts, open the management UI:

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

### Non-Streaming

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

## Features

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

Full config reference: [CLIProxyAPI Documentation](https://help.router-for.me/hands-on/tutorial-0.html)

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

Based on [CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI) v6.8.21 by [router-for-me](https://github.com/router-for-me).
Full documentation: [https://help.router-for.me/](https://help.router-for.me/)
