---
name: dekisugi-email-delivery
description: Email delivery workflow for the DEKISUGI app. Use when changing email configuration, Gmail SMTP, Resend sending, MIME assembly, inline QR attachments, email templates, preview/dry-run behavior, cron queueing, rolling 24-hour quota, or provider validation.
---

# DEKISUGI Email Delivery

Use this for payment email preview, send, provider, and cron work.

## Safety Rule

Never send real email, run a real cron batch, or persist real credentials unless the user explicitly asks for that exact action. Prefer preview and dry-run paths during validation.

## Load Contract

Read `references/email-contract.md` when changing templates, provider logic, cron state, quota, MIME assembly, or local config files.

## Core Invariants

- `email_config.local.json` stores local provider config with mode `0600`.
- Legacy `resend_config.local.json` remains a fallback.
- Gmail is the default provider unless legacy Resend-only config is detected.
- Gmail sends through STARTTLS to `smtp.gmail.com:587`.
- Payment-due email uses inline QR via CID attachment for real sends.
- Preview uses data URL images unless explicitly rendering CID mode.
- Manual sends count against the rolling 24-hour cron quota.
- Transient provider errors stop the current batch and keep remaining cron jobs queued.

## Workflow

1. Identify whether the change is config, rendering, provider send, batch send, or cron.
2. Add a focused test around the public seam:
   - `normalizeEmailConfig`
   - `validateEmailConfigForSend`
   - `renderPaymentEmail`
   - `buildGmailMessage`
   - `sendPaymentEmailRows`
   - `runEmailCronBatch` or cron helpers
3. Keep provider-specific logic behind `sendRenderedEmail`.
4. Preserve masking behavior in public config responses.
5. Keep HTML email table layout simple and email-client-safe: inline CSS, no external CSS.

## Verification

Run:

```sh
go test ./...
```

For manual endpoint checks, use preview or dry-run:

- `POST /api/v1/email/preview`
- `POST /api/v1/email/send` with `"dryRun": true`
- `GET /api/v1/email/cron`

Do not call `POST /api/v1/email/cron/run` unless the user explicitly asks to send the queued batch.
