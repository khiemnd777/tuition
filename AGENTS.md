# Agent Instructions

## Project

DEKISUGI QR Generating System is a small Go app for school payment workflows. It imports payment rows, generates VietQR payloads/PNG images, previews and sends payment emails, and can queue email batches through a local cron state file.

## Runbook

- Install/update modules: `go mod tidy`
- Test everything: `go test ./...`
- Run locally: `go run .`
- Run on another port: `PORT=18081 go run .`
- Open UI: `http://localhost:18080`

## Repo Map

- `main.go`: routes, payment rows, CSV import, QR item assembly.
- `vietqr_standard.go`: VietQR/NAPAS TLV payload and CRC.
- `email.go`: email config, preview, templates, Resend send path.
- `email_gmail.go`: Gmail SMTP and MIME assembly.
- `email_cron.go`: email queue, scheduler, rolling quota.
- `web/`: embedded vanilla HTML/CSS/JS UI.
- `samples/students.csv`: sample import input.
- `main_test.go`: current regression suite.

## Safety

- Do not send real email or run real cron batches unless the user explicitly asks for that exact action.
- Use preview or dry-run for email validation by default.
- Do not commit or print real secrets from `email_config.local.json`, `resend_config.local.json`, or `email_cron.local.json`.
- Preserve `PaymentItems` total overriding raw `Amount` unless the user asks for a contract change.
- Preserve VietQR TLV/CRC behavior with exact tests when touching `vietqr_standard.go`.
- Production UI rule: upsert/detail input workflows must use the app dialog/popup components, not inline panels.
- Do not use native browser `window.alert`, `window.confirm`, or `window.prompt`; use the app dialog/confirm component instead.

## Planning and Review

- For every user request, make a plan for user review before implementation, validation, or file edits.
- Wait for explicit user approval before executing the reviewed plan, unless the user clearly asks only for investigation or status reporting.
- For UI-related requests, including UI changes, UI proposals, layout adjustments, visual design, interaction design, or frontend copy changes, include a `Mock Up as Text` section in the plan so the user can review the intended screen structure before implementation.
- Do not run browser verification automatically. For UI changes, only run browser/browser-verify when the user explicitly asks for it or explicitly approves that validation step in the reviewed plan.
- Keep plans concrete and scoped: include affected files or modules, implementation steps, validation steps, and any assumptions or blockers.

## Agent Skills

Repo-scoped skills live under `.agents/skills` and follow the OpenAI/Codex `SKILL.md` format.

- `$dekisugi-change-workflow`: default workflow for scoped repo changes.
- `$dekisugi-vietqr-payments`: QR payload, payment rows, CSV import, fee totals, PNG output.
- `$dekisugi-email-delivery`: email config, templates, Gmail/Resend, cron, quota.
- `$dekisugi-frontend-ui`: embedded vanilla web UI workflow.
- `$dekisugi-landing-page`: public Landing Page, subscription plan display, Login/Enrollment navigation, hover transforms, and `finance_hub_landing`.
- `$dekisugi-debug-loop`: reproduction-first debugging.
- `$dekisugi-tdd-slice`: one-behavior-at-a-time TDD.

## Custom Subagents

Custom agents live under `.codex/agents`.

- `dekisugi_code_mapper`: read-only codepath exploration.
- `dekisugi_qr_payment_engineer`: QR/payment implementation.
- `dekisugi_email_delivery_engineer`: email/cron implementation.
- `dekisugi_frontend_ui_engineer`: UI implementation and browser verification.
- `dekisugi_landing_page_engineer`: public Landing Page and landing Docker surface implementation.
- `dekisugi_repo_reviewer`: correctness/security/test review.
- `dekisugi_openai_docs_researcher`: OpenAI docs lookup for Codex/skill/subagent questions.

## Documentation

- `CONTEXT.md` is the domain glossary. Keep it implementation-free.
- `docs/adr/` records durable architecture decisions.
- `docs/agents/` explains repo-local agent setup assumptions.
- `docs/landing/DESIGN.md` records Landing Page design, interaction, data, and Docker runtime rules.
- `docs/initiatives/production-module-roadmap.md` records the production module roadmap.
- `docs/initiatives/current-state.md` records cross-session initiative progress.

## Initiative Continuation

When the user asks "Đến đâu rồi?", "bắt đầu tiếp công việc", "tiếp tục", "continue", or "status", read `docs/initiatives/current-state.md` first, then `docs/initiatives/production-module-roadmap.md`.

Report the current phase, completed work, next recommended initiative, and blockers. If the user asks to continue, start the next recommended initiative unless they specify a different initiative. Update `docs/initiatives/current-state.md` after every completed initiative or meaningful partial implementation.
