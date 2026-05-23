---
name: abcsun-tdd-slice
description: Vertical-slice test-driven development workflow for the ABC SUN Go app. Use when adding or fixing behavior test-first in QR generation, CSV import, email rendering, Gmail MIME, cron quota, API handlers, or client/server contracts.
---

# ABC SUN TDD Slice

Use one failing behavior test at a time. Do not write a batch of imagined tests before implementation.

## Test Seams

Prefer these public seams:

- QR payload: `generateVietQR`.
- Payment assembly: `buildQRItem`.
- CSV import: `parseCSVRows`.
- HTTP behavior: handlers with `httptest`.
- Email rendering: `renderPaymentEmail`.
- Gmail MIME: `buildGmailMessage`.
- Cron/quota: `normalizeEmailCronState`, `emailCronDue`, `sentLast24hForState`, `addEmailSentToState`.

Avoid tests that only lock private implementation details unless no public seam can express the behavior.

## Loop

For each behavior:

1. Write one test that describes the observable behavior.
2. Run it and confirm it fails for the right reason.
3. Implement the minimum code to pass.
4. Run the focused test.
5. Run `go test ./...`.
6. Refactor only after the suite is green.

## Good Test Shape

- Name the domain behavior, not the helper detail.
- Use realistic payment rows, CSV fields, or email config values.
- Assert exact strings when payload/MIME/CSV contracts matter.
- Assert structured state when cron or quota behavior matters.
- Keep real network calls out of tests.

## Stop Conditions

Pause and clarify before writing more tests if the public API shape, CSV naming, email provider behavior, or operator workflow is ambiguous.
