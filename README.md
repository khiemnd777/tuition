# ABC SUN Finance Hub

Demo nhỏ để kiểm tra flow sinh VietQR theo danh sách học sinh/phụ huynh trước khi mở rộng import Excel/PDF.

## Chạy nhanh

```sh
go mod tidy
go run .
```

Mở `http://localhost:18080`, bấm `Sinh QR`, rồi scan QR trong preview.
Nếu cổng mặc định đang bận, có thể chạy bằng `PORT=18081 go run .`.

## Chạy bằng Docker

Stack Docker local gồm 4 service: `api`, `admin`, `postgres`, và `redis`.

```sh
cp .env.docker.example .env
docker compose up --build
```

Nếu máy đang có service khác ở các cổng mặc định của stack này, override host ports:

```sh
API_PORT=18182 ADMIN_PORT=18183 POSTGRES_PORT=15437 REDIS_PORT=16382 docker compose up --build
```

Mặc định:

- Admin UI: `http://localhost:18181`
- API trực tiếp: `http://localhost:18180`
- PostgreSQL: `localhost:15436`
- Redis: `localhost:16381`

Trong local Docker, Admin container chạy Vite dev server với HMR. Source repo được bind mount vào `/app`, còn `node_modules` nằm trong Docker volume `admin_node_modules`, nên sửa file trong `web/` sẽ tự reload trên trình duyệt ở `http://localhost:18181`. Vite proxy `/api/` sang `api:18080`, nên UI vẫn dùng các URL tương đối như `/api/v1/banks`.

Khi thêm thư viện frontend mới, cài qua Docker để dependency và lockfile khớp với container:

```sh
docker compose run --rm admin npm install ten-thu-vien
docker compose up -d --build admin
```

Admin container cũng chạy `npm install` khi start để đồng bộ `admin_node_modules` theo `package.json`/`package-lock.json`. Docker build production của admin dùng `npm ci` rồi `npm run build`, nên các import từ npm package sẽ được resolve từ `node_modules` trong image:

```sh
docker build -f docker/admin.Dockerfile --target production -t abcsun-qr-admin:prod .
```

Chạy migration trong Docker:

```sh
docker compose run --rm api migrate up
docker compose run --rm api migrate status
```

Kiểm tra cấu hình DB trong container mà không in secret:

```sh
docker compose run --rm api db config
```

Runtime state local của API, như email config JSON đã ignore khỏi git, nằm trong volume `api_data`. PostgreSQL và Redis lần lượt dùng `postgres_data` và `redis_data`.

Dừng stack:

```sh
docker compose down
```

Reset toàn bộ dữ liệu local Docker:

```sh
docker compose down -v
```

## Codex agents/skills local

Repo này có bộ cấu hình Codex/OpenAI local:

- `AGENTS.md`: chỉ dẫn nền cho agent khi làm việc trong repo.
- `.agents/skills/`: bộ skills repo-scoped cho workflow change, VietQR, email/cron, frontend UI, debugging và TDD.
- `.codex/agents/`: custom subagents cho mapping code, implementation theo domain, review và tra cứu OpenAI docs.
- `CONTEXT.md` và `docs/agents/`: glossary domain và quy ước tiêu thụ context.
- `docs/initiatives/production-module-roadmap.md`: roadmap production theo từng initiative/module để có thể triển khai từng ngày.
- `docs/initiatives/current-state.md`: state hiện tại để agent biết đang đến đâu khi bạn hỏi `Đến đâu rồi?` hoặc `bắt đầu tiếp công việc`.

Khi mở Codex trong thư mục này, có thể gọi trực tiếp các skill như `$abcsun-vietqr-payments`, `$abcsun-email-delivery`, hoặc `$abcsun-debug-loop`.

## Nền tảng persistence production

Server vẫn khởi động được khi chưa cấu hình database để phục vụ static UI và endpoint QR PNG public, nhưng Web Admin production và các API bảo vệ yêu cầu PostgreSQL. Các lệnh production persistence dùng PostgreSQL và đọc cấu hình từ biến môi trường, không lưu secret trong repo:

```sh
export ABC_ENV=local
export ABC_DATABASE_URL_LOCAL='postgres://user:password@localhost:5432/abcsun?sslmode=disable'
go run . db config
go run . db ping
go run . migrate status
go run . migrate up
```

`ABC_ENV` hỗ trợ `local`, `staging`, hoặc `production`. URL database được ưu tiên theo môi trường: `ABC_DATABASE_URL_LOCAL`, `ABC_DATABASE_URL_STAGING`, `ABC_DATABASE_URL_PRODUCTION`; sau đó fallback sang `ABC_DATABASE_URL` và `DATABASE_URL`.

Migration SQL nằm trong `migrations/` và được embed vào binary. Migration runner ghi version/checksum vào `schema_migrations`, nên chạy lại `go run . migrate up` sẽ skip migration đã apply và báo lỗi nếu checksum bị lệch. Schema nền tạo `app_users`, `app_roles`, `app_permissions`, bảng nối role/permission, và `audit_logs` append-only. Schema master data production tạo `schools`, `school_years`, `classes`, `students`, `parents`, và `student_parents`; `student_code` là định danh bắt buộc, nên các học sinh trùng tên vẫn được xử lý an toàn. Schema bảng phí theo kỳ tạo `fee_types`, `fee_schedules`, `fee_schedule_items`, và `student_fee_adjustments` để preview tổng phí trước khi sinh invoice. Schema hóa đơn tạo `invoices`, `invoice_items`, `invoice_adjustments`, `invoice_status_history`, và `receipt_documents` để hóa đơn là nguồn payment request production. Schema thanh toán tạo `payment_providers`, `payment_intents`, `provider_events`, `payment_transactions`, `reconciliation_matches`, và `manual_cash_receipts` để ghi nhận tiền thật và đối soát invoice. Schema vận hành tạo `operation_logs` để giữ lỗi webhook, email và background job cho production review. Schema auth tạo `app_auth_sessions`, `app_auth_access_tokens`, `app_auth_refresh_tokens`, và thêm `password_hash` vào `app_users`; browser token chỉ lưu hash trong DB. Migration RBAC bổ sung permission cho email config/send/cron để tất cả API production có permission route-level rõ ràng. Migration school tree thêm `schools`, backfill dữ liệu hiện có vào `ABC_SUN`, và gắn năm học/lớp vào cây `school > school year/cohort > class`. Migration user/RBAC mới thêm `phone` cho `app_users`, cho phép Email hoặc SĐT là định danh bắt buộc, seed role `admin`, `staff`, `accountant`, và seed permission canonical dạng `{module}.{action}`. Migration tenant foundation tạo `tenants` và `tenant_memberships`, gắn toàn bộ `schools` hiện có vào tenant mặc định `ABC_SUN`, và chuyển unique school code sang phạm vi `(tenant_id, code)` để chuẩn bị cho subscription nhiều trường. Migration tenant-aware auth/RBAC gắn `app_auth_sessions` vào active tenant, tạo `tenant_user_roles`, và backfill role assignment từ `app_user_roles` vào tenant mặc định. Migration tenant data isolation thêm `tenant_id` cho students, parents, notification campaigns, audit logs, operation logs và chuyển unique student/parent/campaign code/email sang phạm vi tenant. Migration tenant onboarding/switching seed quyền tenant, thêm index membership cho switcher, và bật tạo tenant/school khởi tạo từ Web Admin.

Quy ước production hiện tại:

- Primary key dùng UUID do PostgreSQL sinh bằng `gen_random_uuid()`.
- Tenant mặc định hiện là `ABC_SUN`; auth session tự chọn active tenant đầu tiên, route-level RBAC đọc role theo tenant, và các API school-owned production data lọc theo active tenant. Web Admin có tenant switcher, tenant onboarding, và tenant subscription admin cho operator có quyền; phase thanh toán đã triển khai tenant-scoped provider credentials và webhook ownership: provider list sẽ trả `webhookPath` theo tenant, `payOS` ưu tiên credential trong `payment_providers.config`, và write workflow production sẽ bị chặn khi subscription tenant không còn `active` hoặc `trial`. Subscription Phase 8 bổ sung entitlement/usage metering cho `schools`, `operators`, `students`, và `monthly_notifications`; quota được enforce ở các workflow create/import/send thật và tenant panel hiển thị usage theo plan hiện tại.
- Bảng vận hành có `created_at`, `updated_at`, `created_by_user_id`, `updated_by_user_id` khi cần audit actor.
- Phiếu thu tiền mặt và điều chỉnh phí cần lý do; app ghi thêm audit log bất biến cho các thay đổi này.
- Web Admin dùng cookie HttpOnly cho access/refresh token. Access token mặc định sống 15 phút, refresh token mặc định sống 7 ngày và được rotate sau mỗi lần refresh.
- Secret chỉ đi qua environment hoặc secret manager, không commit vào file tracked.
- Runbook backup/restore: `docs/runbooks/backup-restore.md`.
- Runbook deployment, incident, staging smoke test và readiness checklist: `docs/runbooks/production-operations.md`.

Bootstrap admin đầu tiên sau khi chạy migration auth:

```sh
export ABC_AUTH_BOOTSTRAP_EMAIL='admin@example.edu.vn'
export ABC_AUTH_BOOTSTRAP_PHONE='0901234567'
export ABC_AUTH_BOOTSTRAP_PASSWORD='change-this-long-password'
export ABC_AUTH_BOOTSTRAP_DISPLAY_NAME='ABC SUN Admin'
```

Nếu `app_users` chưa có user nào, màn đăng nhập sẽ đổi sang form tạo Admin đầu tiên tại URL Web Admin. Form này yêu cầu password và ít nhất một trong hai trường Email hoặc SĐT; Admin được gán role `admin`. Bootstrap qua biến môi trường vẫn được hỗ trợ: khi login với Email/SĐT bootstrap, app tạo/cập nhật user active và gán role `admin`. Có thể điều chỉnh TTL bằng `ABC_AUTH_ACCESS_TTL` và `ABC_AUTH_REFRESH_TTL` theo định dạng Go duration, ví dụ `15m`, `168h`. Ở production, cookie tự bật `Secure`; local HTTP có thể để mặc định không secure hoặc cấu hình bằng `ABC_AUTH_COOKIE_SECURE`.

Role production mặc định:

- `admin` / `Admin - Quản trị viên`: toàn quyền.
- `staff` / `Staff - Nhân sự`: quản lý học sinh, cây trường/lớp và thông báo.
- `accountant` / `Accountant - Kế toán`: quản lý bảng phí, hóa đơn, thanh toán, đối soát, báo cáo, email/cron và xem audit log.

Permission mới dùng dạng `{module}.{action}`, ví dụ `user.view`, `student.update`, `fee.view`, `invoice.create`, `payment.reconcile`, `report.export`. Các permission cũ vẫn được map tương thích trong code để không làm gãy user/role đã migrate từ phiên bản trước.

Endpoint PNG test nhanh:

```txt
http://localhost:18080/api/v1/qr.png?bankBin=970441&account=625704060370690&amount=0&billNumber=TESTVIB01&note=TEST%20VIB
```

## Excel/CSV mẫu

File mẫu: `samples/students.csv`

Header được hỗ trợ:

```csv
student_name,parent_name,class_name,bank_bin,bank_account,email,amount,bill_number,note
```

`bank_bin` là mã BIN Napas của ngân hàng. VietQR cần BIN + số tài khoản; chỉ có tên ngân hàng hoặc số tài khoản riêng lẻ là chưa đủ để sinh mã chuyển khoản.

`bill_number` map vào Additional Data `62-01`; `note` map vào Purpose of Transaction `62-08`. Nếu không nhập `bill_number`, app tự sinh mã ổn định từ thông tin dòng.

Nếu dùng các khoản phí động, có thể gửi `paymentItems` qua JSON hoặc thêm các cột CSV sau; tổng các cột này sẽ được dùng làm `amount` trong QR và được render thành từng dòng trong email:

```csv
tuition_april,shuttle_april,tuition_may,health_insurance,uniform_fee,international_material,previous_fees
```

Khi file Excel/CSV dùng tên cột riêng, dùng bước `Fields Mapping` trong UI hoặc gửi multipart field `mapping` dạng JSON vào endpoint import. Ví dụ:

```json
{
  "Họ và tên": "student",
  "Phụ huynh": "parent",
  "Tổng phí": "amount"
}
```

Trong UI, workflow được nhóm theo tác vụ production: `Tổng quan`, `Thiết lập`, `Học phí`, `Thu tiền`, `Liên lạc`, và `Báo cáo & vận hành`. Top bar có breadcrumb và bối cảnh trường/năm học/kỳ thu/tháng để đồng bộ nhanh với filter của màn hình đang mở. Mục `Tổng quan` hiển thị việc cần xử lý, bước nhanh theo quyền hiện tại, và Data Quality & Readiness Center để gom blocking/warning/info issue trước khi lập phí, sinh hóa đơn, gửi thông báo hoặc đối soát. Readiness có filter theo severity/loại issue và nút mở nhanh sang học sinh, bảng phí, hóa đơn, thông báo, đối soát, email/cron hoặc operation logs. `Công cụ QR/import` là công cụ phụ cho record thanh toán legacy, không phải nguồn dữ liệu production chính. Màn hình `Học sinh & phụ huynh` có cây trường để quản lý trường, năm học/cohort, khối, lớp, sĩ số, bảng phí và điều chỉnh theo lớp. Cây trường là relationship workspace: mỗi node hiển thị readiness theo kỳ/tháng đang chọn, gồm người nhận billing, bảng phí active, hóa đơn, công nợ đang mở, roster học sinh, và quick action sang danh sách học sinh, thiết lập bảng phí hoặc sinh hóa đơn. Màn hình `Bảng phí` quản lý bảng phí theo kỳ production và vẫn giữ `Template thanh toán tạm` cho record thanh toán legacy; khi bấm `Thêm dòng`, app clone template hiện tại vào dòng mới. Nút `Áp dụng cho tất cả dòng` dùng để đồng bộ template vào các dòng đang có. Các tab production và admin dùng PostgreSQL đã cấu hình.

## Excel/CSV master data

File mẫu: `samples/master_data.csv`

Header production master data:

```csv
student_code,student_name,school,school_year,grade,class_name,parent_name,parent_email,parent_phone,relationship,parent_primary,parent_active,receives_billing_email
```

Quy tắc import:

- `student_code` là bắt buộc và hiện vẫn unique toàn hệ thống để giữ ổn định invoice/payment hiện có.
- `student_name` không được dùng làm định danh; hai học sinh trùng tên phải có `student_code` khác nhau.
- `school` là tùy chọn; nếu bỏ trống app dùng trường mặc định `ABC_SUN`.
- `school_year` và `class_name` là bắt buộc. `grade` có thể bỏ trống nếu app suy ra được từ `class_name`, ví dụ `3.02` -> `3`.
- Một học sinh có thể có nhiều phụ huynh; mỗi học sinh chỉ có một phụ huynh chính đang active.
- `parent_email` được chuẩn hóa lowercase. `parent_phone` được chuẩn hóa về dạng số/SĐT gọn. Nếu `receives_billing_email=true` thì `parent_email` phải có giá trị.
- `relationship` là nhãn quan hệ như `mother`, `father`, `guardian`, `grandparent`, hoặc `other`; nếu bỏ trống app dùng `guardian`.
- Import preview sẽ báo conflict nếu CSV hoặc database hiện có mâu thuẫn; apply import không tự overwrite dữ liệu khác biệt.
- UI sẽ scan header Excel/CSV trước và cho map cột, ví dụ `Họ và tên` -> `Họ, tên`, `Phụ huynh` -> `Tên ba mẹ`.
- Ngoài import batch, tab `Học sinh` có relationship workspace cho học sinh, phụ huynh, người nhận billing, cảnh báo contact, sibling dùng chung phụ huynh, và hóa đơn cần xử lý. Form tạo/sửa thủ công dùng app dialog, upsert theo `studentCode`, chọn lớp hiện có, cập nhật quan hệ phụ huynh/link nhận billing, không xóa phụ huynh cũ nếu không được đánh dấu inactive hoặc bỏ nhận billing. Nếu học sinh đã có invoice hoặc điều chỉnh phí đang active, app chặn đổi lớp để tránh lệch dữ liệu đã phát hành.

## Bảng phí theo kỳ

Bảng phí theo kỳ dùng master data production để thiết lập các khoản phải thu mặc định theo trường, năm học, khối hoặc lớp trước khi sinh invoice. Fee type mặc định gồm `tuition`, `lunch`, `shuttle`, `uniform`, `insurance`, `materials`, `previous_fees`, và `custom`; mỗi loại có nhãn tiếng Việt và tiếng Anh.

Mỗi bảng phí có `periodCode`, `month` tùy chọn, trạng thái `draft` hoặc `active`, các khoản phí mặc định, và điều chỉnh theo học sinh. Điều chỉnh hỗ trợ:

- `discount`: giảm trừ một số tiền.
- `surcharge`: phụ thu.
- `waiver`: miễn giảm; nếu `amount=0` và có `fee_type_code`, app miễn toàn bộ khoản phí mặc định đó.
- `carry_over`: chuyển phí kỳ trước sang kỳ hiện tại.

UI có bảng điều chỉnh theo học sinh cho mã học sinh, loại điều chỉnh, khoản phí, số tiền, và lý do bắt buộc. Ô CSV vẫn hỗ trợ paste nhanh:

```csv
student_code,adjustment_type,fee_type_code,amount,reason
S001,discount,tuition,500000,Ưu đãi anh chị em
S002,waiver,shuttle,0,Không sử dụng xe
```

Preview bảng phí trả tổng mặc định, tổng điều chỉnh, tổng phải thu, và trạng thái người nhận phí cho từng học sinh. Preview báo scope rỗng, class/khối không có học sinh, số tiền không hợp lệ, điều chỉnh thiếu lý do, và học sinh chưa có billing recipient hợp lệ. Danh sách bảng phí đã lưu hiển thị scope, kỳ thu, số khoản, số học sinh, tổng preview, actor/time cập nhật, và có quick action sang preview/generation hóa đơn. Khi có điều chỉnh, payment items dùng cho QR/payment data được gom thành `Tổng phí sau điều chỉnh` để giữ invariant tổng `PaymentItems` quyết định `Amount`.

## Hóa đơn và PDF receipt

Tab `Hóa đơn` là workbench sinh payment request production từ một bảng phí đã lưu. Luồng chính đi theo các bước scope, preview, issues, generate: chọn bảng phí, xem từng dòng hóa đơn dự kiến, xử lý blocking/warning, rồi xác nhận ghi DB. Mỗi học sinh trong phạm vi bảng phí có một hóa đơn active theo `feeScheduleId`; chạy lại cùng bảng phí sẽ trả hóa đơn hiện có để giữ idempotent. Nếu chọn regenerate, app chỉ cập nhật lại hóa đơn chưa có `paid_amount`; hóa đơn đã có tiền thu sẽ bị chặn trong preview.

Preview hóa đơn hiển thị mã hóa đơn, học sinh, lớp, số line item, tổng tiền, trạng thái billing recipient, idempotency state và issue state. Danh sách hóa đơn đã sinh hiển thị tổng tiền, đã thu, còn thiếu, trạng thái, số lần gửi thông báo, QR/PDF readiness và các thao tác QR, PDF, payment intent, notification handoff, export CSV. Panel chi tiết đọc snapshot bất biến từ API detail, gồm line items, adjustments, payment intent/match count, notification sent count và status history.

Hóa đơn snapshot các dữ liệu sau tại thời điểm sinh:

- Mã hóa đơn, học sinh, lớp, năm học, kỳ thu, tháng.
- Dòng phí và điều chỉnh theo học sinh.
- Tổng mặc định, tổng điều chỉnh, tổng phải thu.
- BIN ngân hàng và tài khoản thu.
- `qrBillNumber`, chính là mã hóa đơn dùng cho VietQR `BillNumber`.
- `qrNote` dùng cho VietQR Purpose.

PDF receipt được render từ dữ liệu hóa đơn đã lưu, không đọc lại bảng phí. PDF chứa tên trường, học sinh, lớp, kỳ thu, line items, tổng tiền, trạng thái, thời điểm phát hành, thông tin VietQR, và QR thanh toán.

## Thanh toán và đối soát

Tab `Đối soát` ghi nhận payment intent, webhook provider, giao dịch vào, match đối soát và phiếu thu tiền mặt trên database production. Workbench đối soát có bước Scope, Invoice ledger, Transactions, Match detail và Manual review; filter theo trường, năm học, khối, lớp, kỳ thu, provider, invoice status và transaction status. Summary hiển thị phải thu, đã thu, còn thiếu, collection rate, unpaid, partial, paid, overpaid, manual-review và unmatched để operator thấy ngay hàng đợi cần xử lý. Provider mặc định:

- `manual_vietqr`: dùng QR hóa đơn hiện tại, đối soát thủ công hoặc tiền mặt.
- `sepay`: nhận webhook biến động số dư, lưu raw payload trước rồi parse transaction theo `id`, `transferAmount`, `accountNumber`, `content`, `referenceCode`.
- `payos`: tạo payment link qua API merchant khi có cấu hình env và nhận webhook payment link.

Webhook provider được xử lý idempotent bằng unique key trên `provider_events` và `payment_transactions`; provider retry hoặc gửi lại cùng giao dịch sẽ không tạo thêm khoản thu. Match tự động dựa trên mã hóa đơn hoặc provider reference trong nội dung/reference, số tài khoản thu, và số tiền. Transaction và invoice detail hiển thị match type, match status, score, applied amount, reason, provider reference và collection account để giải thích vì sao giao dịch được hoặc chưa được match. Sau mỗi match, app tính lại `paid_amount` từ `reconciliation_matches` active và cập nhật status invoice thành `unpaid`, `partial`, `paid`, hoặc `overpaid`.

Ghi nhận tiền mặt yêu cầu người thu tiền, số tiền, lý do và mã phiếu thu. Giao dịch tiền mặt cũng đi qua ledger `payment_transactions`, tạo `manual_cash_receipts`, và có match `cash` để audit được như giao dịch webhook. Manual review queue gom invoice partial/overpaid/manual-review và transaction unmatched/manual-review để các trường hợp lệch tiền hoặc thiếu reference dễ tìm.

payOS dùng env sau, không commit secret thật:

```sh
export ABC_PAYOS_CLIENT_ID='...'
export ABC_PAYOS_API_KEY='...'
export ABC_PAYOS_CHECKSUM_KEY='...'
export ABC_PAYOS_RETURN_URL='https://example.edu.vn/payment-return'
export ABC_PAYOS_CANCEL_URL='https://example.edu.vn/payment-cancel'
```

## Notification campaigns

Tab `Thông báo` chuyển luồng gửi email từ danh sách payment row tạm sang campaign dựa trên invoice production. Workbench đi theo Target, Recipients, Email preview, Send/logs, và Cron: campaign chọn invoice theo năm học, kỳ thu, khối/lớp, trạng thái invoice và hạn thanh toán; recipient lấy từ parent contact đang active và có `receives_billing_email=true`.

Template mặc định:

- `first_notice`: thông báo thanh toán lần đầu cho invoice `unpaid`.
- `reminder`: nhắc thanh toán, chỉ cho invoice `unpaid` hoặc `partial`; backend chặn chọn nhầm `paid`.
- `payment_confirmation`: xác nhận đã thanh toán cho invoice `paid`, dùng template `payment_paid` và không đính kèm QR.

Dry-run/preview trả đúng danh sách recipient, số invoice, tổng phải thu, số tiền còn phải thu, QR readiness, trạng thái gửi trước, lỗi và retry eligibility trước khi gửi. Operator có thể chọn từng recipient để render subject/HTML email theo template và invoice trước khi gửi thật. Khi gửi thật, app dùng lại email provider hiện tại, render email học phí với QR inline CID, ghi `notification_logs` theo campaign/template/invoice/recipient/provider/message-id/error/timestamp, và bỏ qua recipient đã gửi trong cùng campaign/template/invoice/email nếu không bật resend rõ ràng. Retry selected yêu cầu campaign đã lưu, recipient được chọn, confirm dialog, `confirmSend=true`, và `forceResend=true`. Email gửi từ campaign cũng tính vào quota rolling 24 giờ giống gửi thủ công và cron.

Khi đối soát webhook hoặc phiếu thu tiền mặt làm invoice chuyển từ `unpaid`/`partial` sang `paid`, app tự gửi email `payment_confirmation` cho billing recipients sau khi transaction thanh toán đã commit. Nếu email config/quota/provider lỗi, payment vẫn giữ thành công và lỗi được ghi vào notification/operation log để operator retry thủ công từ invoice paid.

Panel Cron trong tab `Thông báo` chỉ đọc trạng thái scheduler: enabled, send time, daily limit, queued, sent, errors, sent 24h, next/last run, và recent results. Cấu hình cron vẫn dùng app dialog chung; chạy cron thật chỉ nằm trong tab Email & Cron và luôn cần xác nhận.

## Web Admin

Mục `Tổng quan` tổng hợp công nợ production từ hóa đơn và ledger thanh toán: tổng cần thu, đã thu, còn thiếu, tỷ lệ thu, số học sinh unpaid, partial, overpaid/manual review, giao dịch chưa match và top lớp còn phải thu. Dashboard cũng hiển thị work queue và quick actions theo permission để mở nhanh học sinh, bảng phí, hóa đơn, thông báo, đối soát hoặc công cụ QR/import legacy. Bộ lọc hỗ trợ năm học, khối, lớp, kỳ thu, tháng và trạng thái invoice.

Tab `Báo cáo` dùng cùng bộ lọc để xem tổng hợp theo lớp, chi tiết hóa đơn, giao dịch thanh toán theo provider, và export CSV cho lớp, hóa đơn, giao dịch. Báo cáo giao dịch giữ match explanation gồm match type/status/score, amount applied và reason để đối chiếu provider/reference/account. Tab `Vận hành` đọc `operation_logs` và `audit_logs` để kiểm tra lỗi webhook/email/cron/background job và audit thay đổi tiền/phí; command center có summary lỗi, filter theo operation/status/action/entity type, detail panel metadata đã redacted secret/raw payload, và drilldown sang workflow liên quan khi có entity id. Tab `Người dùng` đọc bảng `app_users`, `app_roles`, `app_permissions`, cho phép tạo/cập nhật user, đặt password tùy chọn, và gán role qua API. Web Admin yêu cầu đăng nhập bằng access/refresh token trước khi gọi API production; backend kiểm tra permission từ role của user đăng nhập cho từng route và trả `403` khi thiếu quyền. UI ẩn menu/action không đủ quyền nhưng không thay thế enforcement backend.

## API

- `GET /api/v1/example`: trả về một VietQR mẫu kèm PNG data URL.
- `POST /api/v1/auth/login`: đăng nhập bằng email/password, set access và refresh cookie HttpOnly.
- `POST /api/v1/auth/refresh`: rotate refresh token và cấp access token mới.
- `POST /api/v1/auth/logout`: revoke session hiện tại và clear cookie.
- `GET /api/v1/auth/session`: trả user/role/permission hiện tại nếu access token còn hiệu lực.
- `GET /api/v1/qr.png`: trả về ảnh PNG để scan trực tiếp.
- `GET /api/v1/banks`: danh sách ngân hàng từ package VietQR.
- `POST /api/v1/import/fields?target=payments|master_data`: scan header Excel/CSV, trả fields và suggested mapping.
- `POST /api/v1/import/csv`: parse Excel/CSV thành rows, hỗ trợ multipart field `mapping`.
- `GET /api/v1/master-data/options`: danh sách trường/năm học/lớp production cho bộ lọc UI.
- `GET /api/v1/master-data/students`: danh sách học sinh production, hỗ trợ `schoolId`, `schoolYearId`, `schoolYear`, `classId`, `grade`, `q`.
- `POST /api/v1/master-data/import/csv?apply=false`: preview import Excel/CSV master data và trả conflict report, hỗ trợ multipart field `mapping`.
- `POST /api/v1/master-data/import/csv?apply=true`: áp dụng import Excel/CSV master data nếu không có conflict.
- `POST /api/v1/master-data/students/save`: tạo/cập nhật thủ công một học sinh, lớp hiện có, và các liên hệ phụ huynh.
- `GET /api/v1/school-tree`: cây `school > school year/cohort > grade > class`, kèm sĩ số, bảng phí và điều chỉnh.
- `POST /api/v1/school-tree/schools/save`: tạo/cập nhật trường.
- `POST /api/v1/school-tree/school-years/save`: tạo/cập nhật năm học/cohort trong một trường.
- `POST /api/v1/school-tree/classes/save`: tạo/cập nhật lớp trong một năm học/cohort.
- `GET /api/v1/fee-schedules/options`: danh sách trường, năm học/lớp và fee type cho bảng phí theo kỳ.
- `GET /api/v1/fee-schedules`: danh sách bảng phí đã lưu, hỗ trợ `schoolId`, `schoolYearId`, `classId`, `grade`, `periodCode`, `month`, `status`.
- `POST /api/v1/fee-schedules/preview`: preview bảng phí theo kỳ trước khi sinh invoice.
- `POST /api/v1/fee-schedules/save`: lưu bảng phí theo kỳ và trả preview mới.
- `GET /api/v1/invoices/options`: danh sách bảng phí, năm học, lớp cho tab hóa đơn.
- `GET /api/v1/invoices`: danh sách hóa đơn, hỗ trợ `schoolYearId`, `classId`, `grade`, `periodCode`, `status`.
- `POST /api/v1/invoices/preview`: preview hóa đơn từ `feeScheduleId` trước khi ghi DB.
- `POST /api/v1/invoices/generate`: sinh hóa đơn idempotent từ bảng phí đã lưu.
- `GET /api/v1/invoices/detail?id=...`: trả snapshot chi tiết, trạng thái QR/PDF, payment/notification counts và status history cho một hóa đơn.
- `GET /api/v1/invoices/payment?id=...`: trả payment row đã sinh QR cho một hóa đơn.
- `GET /api/v1/invoices/pdf?id=...`: trả PDF receipt cho một hóa đơn.
- `GET /api/v1/payments/providers`: danh sách provider thanh toán và trạng thái cấu hình.
- `POST /api/v1/payments/intents`: tạo payment intent cho invoice qua `manual_vietqr`, `sepay`, hoặc `payos`.
- `GET /api/v1/payments/transactions`: danh sách giao dịch ledger, hỗ trợ `provider`, `status`, `limit`.
- `GET /api/v1/payments/reconciliation`: dữ liệu tab đối soát gồm provider, master-data filters, summary, invoices, transactions, intents và reconciliation matches theo invoice.
- `POST /api/v1/payments/webhooks/{provider}`: nhận webhook `sepay` hoặc `payos` theo tenant mặc định `ABC_SUN` (legacy path, vẫn hỗ trợ), lưu raw event, parse transaction và đối soát.
- `POST /api/v1/payments/webhooks/{tenantCode}/{provider}`: nhận webhook tenant-scoped theo tenant code, lưu raw event, parse transaction và đối soát.
- `POST /api/v1/payments/cash-receipts`: ghi nhận phiếu thu tiền mặt vào ledger và invoice; yêu cầu lý do audit qua `reason`.
- `GET /api/v1/subscriptions/plans`: danh sách subscription plan cho tenant admin.
- `POST /api/v1/tenants/subscription/save`: cập nhật plan/status/date của subscription tenant đang active; yêu cầu quyền `subscription.update`.
- `GET /api/v1/notifications/options`: danh sách template, campaign, năm học và lớp cho tab thông báo.
- `GET /api/v1/notifications/templates`: danh sách notification template/version.
- `GET /api/v1/notifications/campaigns`: danh sách campaign đã lưu.
- `POST /api/v1/notifications/campaigns/preview`: dry-run target invoice/recipient trước khi lưu hoặc gửi.
- `POST /api/v1/notifications/campaigns/email-preview`: render subject/HTML cho template và recipient/invoice đã chọn, không gửi email thật.
- `POST /api/v1/notifications/campaigns/save`: lưu campaign và snapshot recipient hiện tại.
- `POST /api/v1/notifications/campaigns/send`: gửi campaign; yêu cầu `confirmSend=true` khi gửi thật, hỗ trợ `recipientIds` cho retry selected và `forceResend=true` để gửi lại có chủ đích.
- `POST /api/v1/notifications/paid-confirmation/send`: gửi hoặc gửi lại confirmation cho một invoice `paid`; yêu cầu `confirmSend=true`.
- `GET /api/v1/notifications/logs`: log gửi theo campaign hoặc gần nhất, hỗ trợ `campaignId`, `limit`.
- `GET /api/v1/healthz`: healthcheck public cho Docker/API liveness.
- `GET /api/v1/admin/dashboard`: dashboard công nợ và readiness center, hỗ trợ `schoolId`, `schoolYearId`, `classId`, `grade`, `periodCode`, `month`, `status`.
- `GET /api/v1/admin/reports`: báo cáo theo lớp, hóa đơn và giao dịch thanh toán, hỗ trợ cùng bộ lọc dashboard và `provider`.
- `GET /api/v1/admin/reports/export?dataset=classes|invoices|transactions`: export CSV theo bộ lọc báo cáo.
- `GET /api/v1/admin/audit-logs`: đọc audit log bất biến, hỗ trợ `action`, `entityType`, `entityId`, `limit`; response có summary command-center và metadata đã redacted secret-like keys.
- `GET /api/v1/admin/operation-logs`: đọc operational log, hỗ trợ `source`, `level`, `operation`, `status`, `entityType`, `entityId`, `limit`; response có summary lỗi webhook/email/cron/background job và metadata đã redacted secret-like keys.
- `GET /api/v1/admin/users`: danh sách user, role, permission.
- `GET/POST /api/v1/auth/bootstrap`: kiểm tra/tạo Admin đầu tiên khi `app_users` rỗng.
- `POST /api/v1/admin/users/save`: tạo/cập nhật user bằng Tên, Email, SĐT, Password; yêu cầu `user.create` hoặc `user.update`.
- `POST /api/v1/admin/users/roles`: gán role `admin`, `staff`, `accountant` cho user; yêu cầu permission `user.assign_role`.
- `GET /api/v1/admin/roles`: danh sách role và permission.
- `POST /api/v1/vietqr/batch`: sinh QR theo danh sách rows.
- `GET/POST /api/v1/email/config`: đọc/lưu cấu hình email local.
- `POST /api/v1/email/preview`: render HTML email theo dòng đầu tiên.
- `POST /api/v1/email/send`: dry-run hoặc gửi email hàng loạt qua Gmail SMTP hoặc Resend.
- `GET/POST /api/v1/email/cron`: đọc/lưu queue cron gửi email.
- `POST /api/v1/email/cron/run`: chạy một đợt cron ngay, vẫn tôn trọng giới hạn ngày.

Repo VietQR bạn gửi là fork của `subiz/vietqr`; module Go trong repo vẫn khai báo `github.com/subiz/vietqr`, nên `go.mod` dùng `replace` sang `github.com/khiemnd777/vietqr`.

Generator payload đã được refactor nội bộ theo tài liệu VietQR/NAPAS v1.5.2:

- `00`: Payload Format Indicator = `01`
- `01`: Point of Initiation Method = `11` hoặc `12`
- `38`: NAPAS AID `A000000727`, BNB BIN, Consumer ID, service `QRIBFTTA`
- `53`: VND `704`
- `54`: số tiền, không có dấu phân tách nghìn
- `58`: `VN`
- `62`: `01` Bill Number và `08` Purpose of Transaction
- `63`: CRC-16/CCITT-FALSE, luôn zero-pad đủ 4 ký tự

## Gửi email

Cấu hình email được lưu ở `email_config.local.json` và đã được ignore khỏi git. Nếu còn file cũ `resend_config.local.json`, app vẫn đọc fallback để không mất cấu hình Resend đã lưu.

Các trường cần nhập trong UI:

- `Provider`: mặc định `Gmail SMTP`; vẫn có thể chọn `Resend`.
- `Gmail address`: tài khoản Gmail/Google Workspace dùng để gửi.
- `Gmail app password`: mật khẩu ứng dụng 16 ký tự của Gmail.
- `Resend API key`: chỉ cần khi chọn Resend, key dạng `re_...`.
- `From`: ví dụ `ABC SUN <billing@example.edu.vn>`. Với Gmail, nếu để trống app dùng `Gmail address`.
- `Reply-To`: email nhận phản hồi
- `Public URL`: domain public của app nếu muốn `QR Link` mở được ngoài máy local

Email thông báo thanh toán dùng QR inline qua CID attachment. Với Gmail, app gửi qua `smtp.gmail.com:587` bằng STARTTLS và gửi tuần tự theo nhịp mặc định 10 email / 5 giây. Với Resend, app vẫn gửi từng email qua endpoint `/emails` vì batch endpoint không hỗ trợ inline attachment.

## Cron gửi email

Tab `Email template` có phần cron để queue danh sách thanh toán hiện tại và gửi dần qua provider đang chọn. Trạng thái cron được lưu ở `email_cron.local.json` và đã được ignore khỏi git.

Thiết lập mặc định:

- `Giới hạn/ngày`: 500 email / 24 giờ.
- `Giờ gửi`: 08:00 theo timezone máy chạy server.
- Mỗi ngày cron chỉ chạy một đợt và chỉ gửi tối đa phần giới hạn còn lại trong rolling window 24 giờ.
- Email gửi thủ công từ app cũng được cộng vào bộ đếm 24 giờ của cron.
- Nếu Gmail/SMTP trả lỗi tạm thời như quota/rate limit, batch hiện tại dừng lại và các email còn lại giữ trạng thái chờ để retry ở lần cron sau.
- `Chạy đợt cron` gửi email thật qua provider đang chọn, nên UI sẽ hỏi xác nhận trước.
