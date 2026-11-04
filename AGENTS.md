# AGENTS.md

Dark factory: the human gives a goal once; the Orchestrator delivers a reviewed PR. The human merges.

## Roles

Provider: Anthropic. Newest generation only: the current lineup on platform.claude.com/docs/en/about-claude/models/overview (Fable 5.1, Opus 5.5, Sonnet 5.5, Haiku 5.5, plus Mythos 5.1 on the pricing page). That page lists Fable 5, Opus 5, Sonnet 5 and older as legacy, so they're excluded.
Prices per MTok input / output, from platform.claude.com/docs/en/about-claude/pricing, checked 2026-11-04.
Only models this account can run are ranked:
- Fable 5.1 ($10 / $50) needs usage credits, so the human dropped it (2026-11-04).
- Mythos 5.1 ($10 / $50) is limited availability. `claude --model claude-mythos-5-1` reports no access.

Price rank, highest first: Opus 5.5 > Sonnet 5.5 > Haiku 5.5.

| Role | Model | Effort | Price |
|---|---|---|---|
| Orchestrator/Planner | `claude-sonnet-5-5` | high | $2 / $10 |
| Executor | `claude-sonnet-5-5` | default | $2 / $10 |
| Auditor | `claude-opus-5-5` | high | $4 / $20 |
| Tooler | `claude-haiku-5-5` | default | $0.10 / $0.50 (prompts ≤100k tokens; $0.50 / $2.50 above) |

- Of the models this account can run, Opus 5.5 is the most capable and the priciest, so it audits.
- Sonnet 5.5 is the best coder that isn't the priciest, so it executes.
- Haiku 5.5 is the cheapest, so it runs Tooler tasks.
- Don't reuse a model across roles unless there's no alternative. The two non-priciest models are already assigned (Executor, Tooler), so the Orchestrator has to reuse one. It takes Sonnet 5.5, run at high effort, because the role needs high reasoning.
- Ask the human only on a price tie.
- A session plays only the role it was given. The top-level session the human talks to is the Orchestrator; every other session is a subagent with one task.
- **Orchestrator/Planner**: started with `claude --model claude-sonnet-5-5 --effort high`. Plans and never codes.
- **Executor**: writes the tests, then the code.
- **Auditor**: fresh context every time. Returns ACCEPT or REJECT. Never fixes. REJECT if any of these hold:
  - tests echo the code;
  - anything is weakened or skipped;
  - a test is deleted;
  - a protected file changed;
  - a dependency is unflagged;
  - any rule here or in the goal is broken.
- The Auditor also tries to break the build with inputs the tests lack. If a break succeeds, REJECT.
- **Tooler**: setup and CI only. Runs on its own, with no Executor or Auditor. Changes to protected files still need the human's sign-off.

## Subagents (herdr)

- Each subagent sees only its own task. Subagents never spawn subagents.
- Run independent tasks in parallel, one herdr pane each:
  1. `herdr pane split --current --direction right --cwd <worktree path> --no-focus` (pane id is `.result.pane.pane_id`)
  2. `herdr agent start <name> --kind claude --pane <pane_id> -- --model <role model> [--effort high] --disallowedTools Agent`
  3. `herdr agent prompt <name> "<task>" --wait` (or `herdr agent wait <name>`)
  4. `herdr agent read <name> --source recent-unwrapped --lines 200`
- Wait only with herdr's blocking wait. Never poll in a loop.

## Flow

1. `git fetch origin main`
2. `wt switch --create <branch> --base origin/main --no-cd --no-hooks -y --format json`. The worktree path is the `.path` field. Every PR branch starts from freshly fetched `main`.
3. The Executor works in that worktree (follow the Go test rules below).
4. Run `task gate` in the worktree. It must pass locally before any PR.
5. Start a new Auditor: ACCEPT or REJECT.
6. On a gate failure or a REJECT, send it back to the Executor and repeat from step 4. If the same failure repeats after two rounds, stop and escalate to the human.
7. On ACCEPT: `git push -u origin <branch>`, then `gh pr create`. The PR body shows the `task gate` result and a per-file coverage table. No custom tooling.
8. The human merges. Then remove the worktree: `wt remove <branch> --no-hooks --foreground -y`
- Create and remove worktrees only with `wt`, with no hooks. Never work in the primary checkout.
- GitHub only via `gh`. Never push to `main`.

## Go tests

- Each exported function gets one table-driven test that calls `t.Parallel()`.
- Each source file gets one blackbox test file (`package <name>_test`).
- Tests assert the goal's behavior, not the code's output.
- Tests are committed before the code they test.
- Keep the CLI `main` thin.
- Coverage must be at least 90%.
- Never weaken a test.

## Protected (changes need the human's sign-off)

- Taskfile, lint config, CI, Dependabot, AGENTS.md, `spec/target.html`, repo settings and rulesets.
- Guardrails and these rules hold unless the human explicitly allows an exception.

## Instruction files

- `CLAUDE.md` is a symlink to `AGENTS.md`. If symlinks fail, use a copy.
