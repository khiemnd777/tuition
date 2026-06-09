# VietQR Payment Contract

## Data Flow

CSV or JSON input becomes `paymentRow`, then `cleanRow`, then `buildQRItem`, then `generateVietQR`, then PNG/data URL output.

Primary public seams:

- `parseCSVRows(input io.Reader) ([]paymentRow, error)`
- `buildQRItem(row paymentRow, size int) qrItem`
- `generateVietQR(req vietQRRequest) (string, error)`
- HTTP handlers in `main.go`

## CSV Fields

Supported canonical fields:

- `student_name`
- `parent_name`
- `class_name`
- `bank_bin`
- `bank_account`
- `email`
- `amount`
- `payment_items`
- `bill_number`
- `note`

Fee columns currently recognized:

- `tuition_april`
- `shuttle_april`
- `tuition_may`
- `health_insurance`
- `uniform_fee`
- `international_material`
- `previous_fees`

When adding a field, update `csvAliases`, tests, and `README.md`.

## TLV Shape

The payload follows this shape:

- `00`: Payload Format Indicator, value `01`
- `01`: Point of Initiation Method, `11` static or `12` dynamic
- `38`: Merchant Account Information for NAPAS
- `38-00`: NAPAS AID, `A000000727`
- `38-01-00`: bank BIN
- `38-01-01`: account number
- `38-02`: service, `QRIBFTTA`
- `53`: VND currency code, `704`
- `54`: amount, only when positive
- `58`: country, `VN`
- `62-01`: bill number when present
- `62-08`: payment purpose when present
- `63`: CRC-16/CCITT-FALSE over payload plus `6304`

## Common Regression Tests

- Exact CRC fixture from VietQR/NAPAS documentation.
- Static payload without amount.
- Dynamic payload with amount, bill number, and purpose.
- Zero-padded CRC suffix.
- Payment-item total overriding raw amount.
- CSV fee columns overriding amount.
