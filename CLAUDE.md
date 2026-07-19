@AGENTS.md
## Data Analysis Review

Any data analysis presented as a conclusion — backtest results, strategy performance claims, Optuna/WFO output, dataset QA findings, or other statistical/quantitative findings — must be adversarially reviewed before being presented to the user. Use one of:

- **An Opus subagent** (`Agent` with `model: "opus"`), explicitly instructed to challenge the methodology — look for overfitting, look-ahead/leakage, cherry-picked windows, confounds, and unsupported causal claims. Not a proofread pass.
- **A Codex adversarial review** (`/codex:adversarial-review`, or the `codex` skill run in review mode) using the strongest available GPT model (currently `gpt-5.6-terra`) at high/xhigh reasoning effort.

Report the adversarial review's findings alongside the analysis itself, not as a separate follow-up step.
