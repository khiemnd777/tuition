# Phát hành Google Sheet mẫu DEKISUGI Email

Đây là thao tác một lần dành cho người phát hành công cụ. Người dùng cuối chỉ cần tạo bản sao, import `DEKISUGI_GMAIL_DATA.json`, gửi thử và gửi thật/đặt lịch.

## 1. Tạo Google Sheet gốc

1. Đăng nhập đúng tài khoản Google sẽ sở hữu template.
2. Mở Google Drive, chọn **New → Google Sheets → Blank spreadsheet**.
3. Đổi tên thành `DEKISUGI Email Sender — Template`.
4. Không nhập dữ liệu phụ huynh thật vào file gốc.

## 2. Gắn Apps Script

1. Trong Sheet, chọn **Extensions → Apps Script**.
2. Mở `Code.gs`, xoá code mẫu và dán toàn bộ nội dung file `Code.gs` trong thư mục này.
3. Bấm dấu **+** cạnh mục **Files**, chọn **HTML**, đặt tên chính xác là `Sidebar`, rồi dán nội dung file `Sidebar.html`.
4. Mở **Project Settings** bằng biểu tượng bánh răng bên trái.
5. Bật **Show "appsscript.json" manifest file in editor**.
6. Quay lại **Editor**, mở `appsscript.json` và thay bằng nội dung file `appsscript.json` trong thư mục này.
7. Bấm biểu tượng **Save project** và chờ trạng thái lưu hoàn tất.

## 3. Cấp quyền và kiểm tra template

1. Trên thanh công cụ Apps Script, mở danh sách tên hàm và chọn `setup`.
2. Bấm **Run**.
3. Khi thấy **Authorization required**, bấm **Review permissions** và chọn đúng tài khoản sở hữu template.
4. Kiểm tra quyền chỉ gồm Sheet hiện tại, gửi email và quản lý trigger của chính script; sau đó bấm **Allow**.
5. Quay lại Sheet và reload trang. Kiểm tra tab `BAT_DAU` đã có hướng dẫn năm bước, sau đó chọn **DEKISUGI Email → Mở bảng điều khiển**.
6. Export một file test từ QR Tool, import vào sidebar và kiểm tra thống kê. Không dùng danh sách phụ huynh thật trong template gốc.
7. Nếu gửi thử, xoá toàn bộ dữ liệu khỏi sheet `EMAILS` sau khi kiểm tra và đảm bảo lịch gửi đang tắt.

## 4. Chia sẻ bằng link “Tạo bản sao”

1. Trong Sheet gốc, bấm **Share**.
2. Đặt **General access → Anyone with the link → Viewer**.
3. Copy link có dạng:

   ```text
   https://docs.google.com/spreadsheets/d/SHEET_ID/edit?usp=sharing
   ```

4. Đổi phần cuối thành `/copy`:

   ```text
   https://docs.google.com/spreadsheets/d/SHEET_ID/copy
   ```

5. Tạo file `qr-tool/.env.local` trên máy build:

   ```text
   VITE_GMAIL_SHEET_TEMPLATE_URL=https://docs.google.com/spreadsheets/d/SHEET_ID/copy
   ```

6. Chạy lại `npm run build`. Không commit `.env.local`; URL template không phải secret nhưng nên được quản lý theo từng môi trường phát hành.

## Checklist an toàn

- File template gốc không chứa email, QR hoặc dữ liệu thanh toán thật.
- Không có trigger `runScheduledBatch` trong template gốc.
- Sidebar không dùng `UrlFetchApp` và không gửi dữ liệu tới DEKISUGI.
- Nút gửi thật và bật lịch chỉ mở sau khi bản sao đã gửi thử thành công.
- Link public cấp quyền **Viewer**, không cấp **Editor**.
