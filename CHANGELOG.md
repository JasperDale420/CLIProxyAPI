# Changelog

## Unreleased

### Fixed
- Removed 5 unreachable duplicate `return` statements in token-count paths across `antigravity_executor.go`, `gemini_cli_executor.go`, `gemini_executor.go`, and `gemini_vertex_executor.go` (flagged by `go vet`). Logic is unchanged — the reachable return already used the correct `[]byte` value.

### Changed
- Synced upstream CLIProxyAPI from `v6.8.51` to `v6.9.6` — 63 merge conflicts resolved across 238 files. Upstream adds weighted provider rotation, Claude max_tokens defaults, Codex capacity error retries, batch auth file upload/delete, auth file name validation (security), and FreeBSD build support.

### Removed
- Removed `gpt-4o-mini` model alias from GLM provider config.
- Removed `deepseek-coder` alias — this model does not exist in DeepSeek's API.

### Changed
- Scheduled maintenance check 2026-03-20: `go build ./...`, `go test ./...`, and `go test -race ./...` all pass with zero errors or data races. All test packages with test files pass. No code-level fixes required.
- Preflight check 2026-03-15: `go build ./...`, `go vet ./...`, and `go test ./...` all pass with no errors. All 18 test packages with test files pass (cached). No fixes required.

### Fixed
- Ignore local cache/editor/runtime artifacts in .gitignore to reduce uncommitted noise.
- Fixed data race in `GetContextWithCancel`: the background goroutine that watches for request context cancellation now captures the cancellable context by parameter instead of closure, preventing a race with the subsequent `context.WithValue` reassignments.
- Fixed data race in `TestPersistConfigAndAuthAsyncInvokePersister`: `stubStore.lastAuthMessage` and `lastAuthPaths` are now protected by a mutex, eliminating an unsynchronized cross-goroutine read/write.
- Fixed data race in `TestScheduleConfigReloadDebounces`: test now reads `w.lastConfigHash` under `clientsMutex.RLock()`, consistent with the protection used by production code.

### Changed
- Synced this repo to upstream CLIProxyAPI `v6.8.51` from `v6.8.21` and refreshed connection handling, including Antigravity transport and cancellation fixes.
- Updated `README.md` and `config.example.yaml` for the v6.8.51 baseline.


## [Unreleased]

### Changed

- Switched the built-in GLM provider from the general Z.ai API endpoint to the Z.ai coding endpoint so gateway traffic matches Coding Plan subscriptions.
- Reduced the advertised GLM model list to the coding-endpoint models exposed by `/models`: `glm-4.5`, `glm-4.5-air`, `glm-4.6`, `glm-4.7`, and `glm-5`.
- Updated the README model reference to document the coding-endpoint GLM routing behavior.

## [1.0.3] - 2026-03-05

### Changed

- Updated the default GLM/z.ai OpenAI-compatible base URL to `https://api.z.ai/api/paas/v4`
- Updated local GLM provider config to use the current z.ai production host instead of the legacy BigModel host

### Fixed

- Added a regression test that guards the built-in iFlow/GLM default API base URL against drifting away from the current z.ai endpoint

## [1.0.2] - 2026-02-18

### Changed

- Port changed from 8001 to 8002 to avoid conflict with Shared-MCP-Server

## [1.0.1] - 2026-02-18

### Added

- Empire-specific README.md with usage guide, model reference, and configuration docs
- Management UI login instructions and management API examples
- Troubleshooting section

## [1.0.0] - 2026-02-18

### Changed

- Replaced custom Python/FastAPI AI-Gateway implementation with CLIProxyAPI (v6.8.21)
- Switched from Python subprocess-based CLI execution to Go-native proxy architecture
- Port initially set to 8001 (changed to 8002 in v1.0.2)

### Added

- Streaming response support
- Function calling / tools support
- Multimodal input support (text and images)
- Multi-account round-robin load balancing
- API key authentication
- Management API and web control panel
- Support for additional providers: Claude Code, Qwen Code, iFlow
- Configurable retry logic with quota-exceeded auto-switching
- OpenAI Responses API compatibility

### Removed

- Custom Python source code (`src/ai_gateway/`)
- Python test suite (`tests/`)
- Python build config (`pyproject.toml`)
- Custom Dockerfile (replaced by CLIProxyAPI's Go-based Dockerfile)

## 2026-02-21
- chore: workspace sync checkpoint and gitignore audit (2026-02-21)
