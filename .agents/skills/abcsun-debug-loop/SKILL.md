---
name: abcsun-debug-loop
description: Reproduction-first debugging workflow for the ABC SUN VietQR app. Use when a QR payload, CSV import, email preview/send, cron batch, Gmail/Resend provider, or browser UI flow is broken, failing, flaky, slow, or producing unexpected output.
---

# ABC SUN Debug Loop

Build a fast, deterministic feedback loop before changing code.

## Phase 1: Feedback Loop

Choose the narrowest loop that reproduces the symptom:

1. `go test ./...` or a focused Go test.
2. `httptest` around a handler.
3. A fixture CSV through `parseCSVRows` or `/api/v1/import/csv`.
4. A direct call to `generateVietQR`, `buildQRItem`, `renderPaymentEmail`, or cron helpers.
5. `curl` against a running local server.
6. Browser reproduction for UI-only issues.

For email issues, use preview or dry-run first. Do not trigger real sends unless explicitly requested.

## Phase 2: Reproduce

Confirm the loop matches the user-visible bug. Capture the exact wrong payload, HTTP response, email HTML/MIME snippet, cron state, browser error, or timing.

## Phase 3: Hypotheses

Write 3-5 ranked, falsifiable hypotheses before editing. Each hypothesis should name the prediction that would prove or disprove it.

Examples:

- If `PaymentItems` total is overriding `Amount`, then a row without items should generate the expected amount.
- If CID mismatch breaks inline QR, then `src="cid:..."` and `Content-ID` will differ in `buildGmailMessage`.
- If cron quota migration is wrong, then `SentToday` without `SendHistory` should miscount only on same local date.

## Phase 4: Instrument

Add targeted probes only where they distinguish hypotheses. Prefix temporary logs with `[DEBUG-abc...]` and remove them before finishing.

## Phase 5: Fix and Lock

Turn the minimal reproduction into a regression test at a public seam, then make the smallest fix. Re-run the original loop and `go test ./...`.

## Cleanup

Before reporting done:

- Remove temporary debug logs and harnesses.
- Preserve useful fixtures only if they become tests.
- State the root cause and the verification signal.
