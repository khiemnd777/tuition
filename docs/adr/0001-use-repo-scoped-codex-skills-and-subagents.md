# ADR 0001: Use repo-scoped Codex skills and custom subagents

## Status

Accepted

## Context

This repo has compact but high-risk workflows: VietQR payload correctness, parent email delivery, cron quota behavior, and a vanilla embedded web UI. Generic agent instructions are not enough to keep future changes consistently scoped and safe.

## Decision

Use repo-scoped Codex skills in `.agents/skills` and custom subagents in `.codex/agents`.

The root `AGENTS.md` provides baseline repo instructions. Domain vocabulary lives in `CONTEXT.md`. Durable agent setup assumptions live under `docs/agents`.

## Consequences

- Codex can discover local skills from `.agents/skills` when launched in this folder.
- Subagents can be spawned by name for focused exploration, implementation, review, or docs lookup.
- Skill and agent behavior is versioned with the local repo folder instead of depending only on user-global configuration.
- Future maintainers should update these artifacts when payment, email, cron, or UI ownership changes.
