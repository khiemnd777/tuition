---
name: dekisugi-change-workflow
description: Repo-specific change workflow for the DEKISUGI VietQR Go app. Use when implementing, reviewing, or validating any non-trivial code change in this repository, especially changes that cross Go handlers, VietQR generation, email/cron behavior, or the vanilla web UI.
---

# DEKISUGI Change Workflow

Use this as the default workflow for changes in this repo. Keep edits narrow and preserve the current simple Go + embedded static-web architecture.

## Orient

Read these first:

- `AGENTS.md` for repo rules and verification commands.
- `CONTEXT.md` for domain vocabulary.
- `README.md` for supported API and operator workflow.
- The smallest owning code path for the request.

File map:

- `main.go`: HTTP routes, payment rows, CSV import, QR item assembly.
- `vietqr_standard.go`: NAPAS/VietQR TLV payload generation and CRC.
- `email.go`: config, templates, preview, Resend sender, email batch flow.
- `email_gmail.go`: Gmail SMTP and MIME assembly.
- `email_cron.go`: local email queue, rolling quota, scheduler.
- `web/*`: embedded vanilla HTML/CSS/JS UI.
- `main_test.go`: regression tests for QR, CSV, email, Gmail MIME, and cron.

## Classify Scope

Choose the narrower companion skill when useful:

- QR payload, CSV import, payment rows, fees: `$dekisugi-vietqr-payments`.
- Email template, provider, SMTP, Resend, cron, quota: `$dekisugi-email-delivery`.
- Browser UI, tab workflow, layout, client API wiring: `$dekisugi-frontend-ui`.
- Bug investigation: `$dekisugi-debug-loop`.
- Test-first feature or fix: `$dekisugi-tdd-slice`.

## Implement

Work through the public behavior first, then the implementation.

1. Identify the handler/function/UI flow that owns the behavior.
2. Add or update a focused test when the change affects QR, CSV, email, cron, validation, or MIME output.
3. Keep data contracts stable unless the user explicitly asks for an API shape change.
4. Update `README.md` only when an operator-facing command, CSV field, endpoint, or workflow changes.
5. Never write real credentials into tracked files. Treat `email_config.local.json`, `resend_config.local.json`, and `email_cron.local.json` as local runtime state.

## Validate

Always run:

```sh
go test ./...
```

Add these checks by scope:

- API or server changes: run `go run .` or `PORT=8081 go run .`, then hit the relevant `/api/v1/...` endpoint.
- UI changes: verify the running app in a browser and check the tab/workflow you changed.
- Email changes: preview or dry-run first; do not send real email or run real cron unless explicitly requested.
- QR changes: assert exact payload, TLV fields, CRC, and PNG/data URL behavior where applicable.

## Report

Finish with changed files, validation performed, and any residual risk. If you skipped a relevant validation, say why.
