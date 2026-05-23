# Email Delivery Contract

## File Map

- `email.go`: config, public config response, preview, batch send, templates, Resend HTTP call.
- `email_gmail.go`: Gmail SMTP, MIME headers, multipart/related + alternative, inline QR attachment.
- `email_cron.go`: local queue, scheduler, rolling quota, cron public state.
- `main_test.go`: email, Gmail MIME, and cron regression tests.

## Provider Contract

Provider values:

- `gmail`
- `resend`

Gmail requires:

- `GmailAddress`
- `GmailAppPassword`

Resend requires:

- `APIKey`
- `From`

Config responses must expose only masked secrets and boolean presence flags.

## Send Contract

`sendPaymentEmailRow` returns statuses:

- `sent`
- `dry_run`
- `skipped`
- `error`

Rows with invalid QR data, invalid email, or template errors do not send.

## Cron Contract

Cron state lives in `email_cron.local.json`.

Important rules:

- Daily limit caps at 500.
- Quota is rolling 24 hours, not calendar-day-only.
- Cron runs once per local date after configured send time.
- Manual sends record into the same quota state.
- Transient errors leave the current job queued and stop the batch.

## High-Risk Areas

- MIME header formatting and CRLF line endings.
- Content-ID matching between HTML `cid:` source and attachment header.
- Secret masking and local config persistence.
- Rolling quota migration from legacy `SentToday`.
- Real provider calls during tests or manual validation.
