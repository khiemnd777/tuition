# Finance Hub Demo Sample Pack

This folder contains Excel workbooks for demonstrating additional imports after the Finance Hub demo database has already been seeded.

## Database Sample Data

`go run . migrate up` automatically applies the versioned demo data migration `0001 finance_hub_demo` after schema migrations. Therefore `make up` on a fresh database creates the full demo tenant without any manual import step.

Seeded workspace:

- Tenant `SUNRISE_DEMO`
- School `SUNRISE`
- Owner `owner.demo@example.com`
- Password `DemoOwner@2026!`
- Owner role `Tenant Owner` with full demo tenant management permissions
- Subscription plans `Free`, `Go`, `Plus`, and `Pro`
- Tenant subscription on `Pro`
- Master data, fee schedules, invoices, VietQR payment intents, payment/reconciliation cases, dry-run notification campaigns, and operation-log review cases

To refresh generated demo operational data in an already-seeded database, run:

```sh
go run . demo seed finance-hub --refresh
```

## Excel Files

- `import_more_students.xlsx`: imports additional students and parent contacts into `Học sinh & phụ huynh`.
- `import_more_payments.xlsx`: imports additional payment rows into `Công cụ QR/import` for QR preview.

These files are not required for initial setup. They exist only to show customers that Excel import can add more data after the demo database is already populated.

## Suggested Demo Flow

1. Run `make up` on a fresh database.
2. Log in as `owner.demo@example.com` with `DemoOwner@2026!`.
3. Review existing seeded records in `Học sinh & phụ huynh`, `Bảng phí`, `Hóa đơn`, `Đối soát`, `Thông báo`, and `Gói & Thanh toán`.
4. Import `import_more_students.xlsx` to show additional master data being added.
5. Import `import_more_payments.xlsx` to show the secondary QR/import workflow.

## Safety

- Do not send real email or run real cron for this sample pack.
- Do not use placeholder bank accounts as real payment instructions.
- Keep webhook payload templates in preview/demo environments only.
