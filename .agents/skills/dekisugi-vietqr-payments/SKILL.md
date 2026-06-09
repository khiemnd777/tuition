---
name: dekisugi-vietqr-payments
description: VietQR payment-row and CSV workflow for the DEKISUGI app. Use when changing VietQR payload generation, NAPAS TLV fields, CRC, bank BIN validation, payment items, amount calculation, bill numbers, notes, CSV import aliases, QR PNG output, or payment API contracts.
---

# DEKISUGI VietQR Payments

Use this for changes to payment rows, CSV import, VietQR payloads, and generated QR output.

## Load Contract

Read `references/vietqr-contract.md` when the change touches TLV payload shape, CSV fields, amount rules, or public payment APIs.

## Core Invariants

- `generateVietQR` is the only place that assembles the EMV/NAPAS TLV payload.
- `cleanRow` canonicalizes input before QR generation.
- `PaymentItems` override `Amount` when present; the total fee list becomes the QR amount.
- `BillNumber` maps to Additional Data `62-01`.
- `Note` maps to Purpose of Transaction `62-08`.
- `BillNumber` and `Note` are ASCII-normalized ANS strings capped at 25 characters.
- Bank BIN must be 6 digits and exist in `vietqr.VNBankM`.
- Account number must be non-empty and at most 19 alphanumeric characters.
- QR is dynamic when `Amount > 0`; static otherwise unless explicitly requested.
- CRC must be CRC-16/CCITT-FALSE, uppercase, and zero-padded to 4 characters.

## Workflow

1. Start from the behavior surface: CSV import, batch endpoint, PNG endpoint, or direct `generateVietQR`.
2. Add a regression test in `main_test.go` before changing payload, CSV, validation, or amount behavior.
3. Keep tests at public seams: `parseCSVRows`, `buildQRItem`, `generateVietQR`, or handlers via `httptest`.
4. Update README CSV/API docs when adding or renaming a CSV column or JSON field.
5. Keep QR payload changes explicit in tests with exact expected substrings or exact payloads.

## Verification

Run:

```sh
go test ./...
```

For endpoint behavior, also exercise:

- `POST /api/v1/import/csv`
- `POST /api/v1/vietqr/batch`
- `GET /api/v1/qr.png`

Do not loosen validation to make malformed rows pass silently unless the user explicitly asks for a tolerant import mode.
