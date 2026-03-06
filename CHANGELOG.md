# Changelog

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
