# Domain Docs

This repo uses a single-context documentation layout:

- `CONTEXT.md` at the repo root is the shared domain glossary.
- `docs/adr/` contains durable architecture decisions.
- `.agents/skills/*` contains Codex skills that should read the glossary when domain vocabulary matters.

Rules for agents:

- Read `CONTEXT.md` before naming new payment, QR, email, or cron concepts.
- Keep `CONTEXT.md` about domain language only. Do not put implementation plans or API specs there.
- Add an ADR only for decisions that are hard to reverse, surprising without context, and the result of a real tradeoff.
