# Agent Instructions

## Project

DEKISUGI QR Tool in `qr-tool/` is the official product. It is a static browser app for importing spreadsheet payment rows, generating VietQR PNG files, preparing email drafts, and exporting user-owned Gmail workflows without a DEKISUGI backend.

The former Go Finance Hub, including `main.go`, `web/`, PostgreSQL modules, subscription/tenant flows, email sender, and cron scheduler, is obsolete and planned for complete removal. Do not use or extend that legacy architecture for new work unless the user explicitly requests a legacy change.

## Official App Runbook

- Install dependencies: `cd qr-tool && npm ci`
- Test: `cd qr-tool && npm test`
- Run locally: `cd qr-tool && npm run dev`
- Build: `cd qr-tool && npm run build`
- Open UI: `http://localhost:5277`

## Repo Map

- `qr-tool/index.html`: official app structure and dialogs.
- `qr-tool/src/main.js`: browser state, rendering, and interactions.
- `qr-tool/src/styles.css`: official app visual system and responsive layout.
- `qr-tool/src/coffee.js`: public coffee-support recipient and VietQR assembly.
- `qr-tool/src/vietqr.js`: VietQR/NAPAS TLV payload and CRC.
- `qr-tool/src/importer.js`: spreadsheet parsing and field mapping.
- `qr-tool/src/email.js`: safe email template and EML assembly.
- `qr-tool/src/exporter.js`: QR, Gmail, and provider exports.
- `qr-tool/tests/`: official regression suite.
- Go files and `web/`: obsolete Finance Hub retained only until its future removal.

## Safety

- Do not send real email or run real cron batches unless the user explicitly asks for that exact action.
- Use preview or dry-run for email validation by default.
- Do not commit or print real secrets from `email_config.local.json`, `resend_config.local.json`, or `email_cron.local.json`.
- Preserve `PaymentItems` total overriding raw `Amount` unless the user asks for a contract change.
- Preserve VietQR TLV/CRC behavior with exact tests when touching `qr-tool/src/vietqr.js` or official QR assembly code.
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
- `$dekisugi-vietqr-payments`: QR payload, payment rows, Excel import, fee totals, PNG output.
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
