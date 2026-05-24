# ABC SUN - QR Generating System

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

Server QR/email mặc định vẫn chạy không cần database. Các lệnh production persistence dùng PostgreSQL và đọc cấu hình từ biến môi trường, không lưu secret trong repo:

```sh
export ABC_ENV=local
export ABC_DATABASE_URL_LOCAL='postgres://user:password@localhost:5432/abcsun?sslmode=disable'
go run . db config
go run . db ping
go run . migrate status
go run . migrate up
```

`ABC_ENV` hỗ trợ `local`, `staging`, hoặc `production`. URL database được ưu tiên theo môi trường: `ABC_DATABASE_URL_LOCAL`, `ABC_DATABASE_URL_STAGING`, `ABC_DATABASE_URL_PRODUCTION`; sau đó fallback sang `ABC_DATABASE_URL` và `DATABASE_URL`.

Migration SQL nằm trong `migrations/` và được embed vào binary. Migration runner ghi version/checksum vào `schema_migrations`, nên chạy lại `go run . migrate up` sẽ skip migration đã apply và báo lỗi nếu checksum bị lệch. Schema nền tạo `app_users`, `app_roles`, `app_permissions`, bảng nối role/permission, và `audit_logs` append-only. Schema master data production tạo `school_years`, `classes`, `students`, `parents`, và `student_parents`; `student_code` là định danh bắt buộc, nên các học sinh trùng tên vẫn được xử lý an toàn. Schema bảng phí theo kỳ tạo `fee_types`, `fee_schedules`, `fee_schedule_items`, và `student_fee_adjustments` để preview tổng phí trước khi sinh invoice. Schema hóa đơn tạo `invoices`, `invoice_items`, `invoice_adjustments`, `invoice_status_history`, và `receipt_documents` để hóa đơn là nguồn payment request production.

Quy ước production hiện tại:

- Primary key dùng UUID do PostgreSQL sinh bằng `gen_random_uuid()`.
- Bảng vận hành có `created_at`, `updated_at`, `created_by_user_id`, `updated_by_user_id` khi cần audit actor.
- Secret chỉ đi qua environment hoặc secret manager, không commit vào file tracked.
- Runbook backup/restore: `docs/runbooks/backup-restore.md`.

Endpoint PNG test nhanh:

```txt
http://localhost:18080/api/v1/qr.png?bankBin=970441&account=625704060370690&amount=0&billNumber=TESTVIB01&note=TEST%20VIB
```

## CSV mẫu

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

Trong UI, workflow được tách thành các tab: `Học sinh`, `Bảng phí`, `Hóa đơn`, `Thanh toán`, và `Email & Cron`. Tab `Bảng phí` quản lý bảng phí theo kỳ production và vẫn giữ `Template thanh toán tạm` cho record thanh toán legacy; khi bấm `Thêm dòng`, app clone template hiện tại vào dòng mới. Nút `Áp dụng cho tất cả dòng` dùng để đồng bộ template vào các dòng đang có. Tab `Học sinh`, `Bảng phí`, và `Hóa đơn` dùng PostgreSQL đã cấu hình.

## CSV master data

File mẫu: `samples/master_data.csv`

Header production master data:

```csv
student_code,student_name,school_year,grade,class_name,parent_name,parent_email,parent_primary,parent_active,receives_billing_email
```

Quy tắc import:

- `student_code` là bắt buộc và unique trong phạm vi trường hiện tại.
- `student_name` không được dùng làm định danh; hai học sinh trùng tên phải có `student_code` khác nhau.
- `school_year` và `class_name` là bắt buộc. `grade` có thể bỏ trống nếu app suy ra được từ `class_name`, ví dụ `3.02` -> `3`.
- Một học sinh có thể có nhiều phụ huynh; mỗi học sinh chỉ có một phụ huynh chính đang active.
- `parent_email` được chuẩn hóa lowercase. Nếu `receives_billing_email=true` thì `parent_email` phải có giá trị.
- Import preview sẽ báo conflict nếu CSV hoặc database hiện có mâu thuẫn; apply import không tự overwrite dữ liệu khác biệt.

## Bảng phí theo kỳ

Bảng phí theo kỳ dùng master data production để thiết lập các khoản phải thu mặc định theo năm học, khối hoặc lớp trước khi sinh invoice. Fee type mặc định gồm `tuition`, `lunch`, `shuttle`, `uniform`, `insurance`, `materials`, `previous_fees`, và `custom`; mỗi loại có nhãn tiếng Việt và tiếng Anh.

Mỗi bảng phí có `periodCode`, `month` tùy chọn, trạng thái `draft` hoặc `active`, các khoản phí mặc định, và điều chỉnh theo học sinh. Điều chỉnh hỗ trợ:

- `discount`: giảm trừ một số tiền.
- `surcharge`: phụ thu.
- `waiver`: miễn giảm; nếu `amount=0` và có `fee_type_code`, app miễn toàn bộ khoản phí mặc định đó.
- `carry_over`: chuyển phí kỳ trước sang kỳ hiện tại.

Ô điều chỉnh trong UI nhận CSV ngắn:

```csv
student_code,adjustment_type,fee_type_code,amount,reason
S001,discount,tuition,500000,Ưu đãi anh chị em
S002,waiver,shuttle,0,Không sử dụng xe
```

Preview bảng phí trả tổng mặc định, tổng điều chỉnh và tổng phải thu cho từng học sinh. Khi có điều chỉnh, payment items dùng cho QR/payment data được gom thành `Tổng phí sau điều chỉnh` để giữ invariant tổng `PaymentItems` quyết định `Amount`.

## Hóa đơn và PDF receipt

Tab `Hóa đơn` sinh payment request production từ một bảng phí đã lưu. Mỗi học sinh trong phạm vi bảng phí có một hóa đơn active theo `feeScheduleId`; chạy lại cùng bảng phí sẽ trả hóa đơn hiện có để giữ idempotent. Nếu chọn regenerate, app chỉ cập nhật lại hóa đơn chưa có `paid_amount`.

Hóa đơn snapshot các dữ liệu sau tại thời điểm sinh:

- Mã hóa đơn, học sinh, lớp, năm học, kỳ thu, tháng.
- Dòng phí và điều chỉnh theo học sinh.
- Tổng mặc định, tổng điều chỉnh, tổng phải thu.
- BIN ngân hàng và tài khoản thu.
- `qrBillNumber`, chính là mã hóa đơn dùng cho VietQR `BillNumber`.
- `qrNote` dùng cho VietQR Purpose.

PDF receipt được render từ dữ liệu hóa đơn đã lưu, không đọc lại bảng phí. PDF chứa tên trường, học sinh, lớp, kỳ thu, line items, tổng tiền, trạng thái, thời điểm phát hành, thông tin VietQR, và QR thanh toán.

## API

- `GET /api/v1/example`: trả về một VietQR mẫu kèm PNG data URL.
- `GET /api/v1/qr.png`: trả về ảnh PNG để scan trực tiếp.
- `GET /api/v1/banks`: danh sách ngân hàng từ package VietQR.
- `POST /api/v1/import/csv`: parse CSV thành rows.
- `GET /api/v1/master-data/options`: danh sách năm học/lớp production cho bộ lọc UI.
- `GET /api/v1/master-data/students`: danh sách học sinh production, hỗ trợ `schoolYearId`, `schoolYear`, `classId`, `grade`, `q`.
- `POST /api/v1/master-data/import/csv?apply=false`: preview import CSV master data và trả conflict report.
- `POST /api/v1/master-data/import/csv?apply=true`: áp dụng import master data nếu không có conflict.
- `GET /api/v1/fee-schedules/options`: danh sách năm học/lớp và fee type cho bảng phí theo kỳ.
- `GET /api/v1/fee-schedules`: danh sách bảng phí đã lưu, hỗ trợ `schoolYearId`, `classId`, `grade`, `status`.
- `POST /api/v1/fee-schedules/preview`: preview bảng phí theo kỳ trước khi sinh invoice.
- `POST /api/v1/fee-schedules/save`: lưu bảng phí theo kỳ và trả preview mới.
- `GET /api/v1/invoices/options`: danh sách bảng phí, năm học, lớp cho tab hóa đơn.
- `GET /api/v1/invoices`: danh sách hóa đơn, hỗ trợ `schoolYearId`, `classId`, `grade`, `periodCode`, `status`.
- `POST /api/v1/invoices/preview`: preview hóa đơn từ `feeScheduleId` trước khi ghi DB.
- `POST /api/v1/invoices/generate`: sinh hóa đơn idempotent từ bảng phí đã lưu.
- `GET /api/v1/invoices/payment?id=...`: trả payment row đã sinh QR cho một hóa đơn.
- `GET /api/v1/invoices/pdf?id=...`: trả PDF receipt cho một hóa đơn.
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
