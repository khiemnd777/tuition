# Demo Sample Data

The Finance Hub demo database is seeded automatically by `migrate up`. After a fresh database reset, `make up` creates the schema and inserts the full fictional tenant sample data into PostgreSQL.

## Default Seeded Data

The default sample tenant is `SUNRISE_DEMO` with owner login:

- Email: `owner.demo@example.com`
- Password: `DemoOwner@2026!`

The seeded database includes tenant setup, school/year/class tree, students, parents, fee schedules, invoices, VietQR payment intents, payment/reconciliation cases, dry-run notification campaigns, operation logs, and subscription plans `Free / Go / Plus / Pro`.

No spreadsheet import is required to get the demo database ready.

## Excel Import Workbooks

Use these `.xlsx` files only to demonstrate adding more data after the database already has sample data:

- `finance_hub_demo/import_more_students.xlsx`: additional students and parent contacts for the master data import workflow.
- `finance_hub_demo/import_more_payments.xlsx`: additional legacy QR/import payment rows for the secondary QR tool.

## Demo Flow

1. Run `make up` after resetting the database.
2. Log in as `owner.demo@example.com` with `DemoOwner@2026!`.
3. Review the seeded Finance Hub data across `Học sinh & phụ huynh`, `Bảng phí`, `Hóa đơn`, `Đối soát`, `Thông báo`, and `Gói & Thanh toán`.
4. Import `finance_hub_demo/import_more_students.xlsx` from `Học sinh & phụ huynh` to show that Excel import adds more records.
5. Import `finance_hub_demo/import_more_payments.xlsx` from `Công cụ QR/import` to show the legacy QR import path.

All data is fictional. Emails use `example.com`, collection accounts are placeholders, and notification campaigns are dry-run safe.
