package main

import (
	"archive/zip"
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDocumentStaticCRCExample(t *testing.T) {
	input := "00020101021138570010A00000072701270006970403011200110123456780208QRIBFTTA53037045802VN6304"
	if got := crc16CCITTFalse(input); got != "F4E5" {
		t.Fatalf("unexpected CRC, want F4E5 got %s", got)
	}
}

func TestStaticIBFTAccountPayload(t *testing.T) {
	got, err := generateVietQR(vietQRRequest{
		BankBIN:       "970423",
		AccountNumber: "0099999999",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "00020101021138540010A00000072701240006970423011000999999990208QRIBFTTA53037045802VN6304CBB4"
	if got != want {
		t.Fatalf("unexpected payload\nwant: %s\n got: %s", want, got)
	}
}

func TestSpecDynamicIBFTAccountPayload(t *testing.T) {
	got, err := generateVietQR(vietQRRequest{
		BankBIN:       "970403",
		AccountNumber: "0011012345678",
		Amount:        180000,
		BillNumber:    "NPS6869",
		Purpose:       "thanh toan don hang",
		Dynamic:       true,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "00020101021238570010A00000072701270006970403011300110123456780208QRIBFTTA530370454061800005802VN62340107NPS68690819thanh toan don hang63042E2E"
	if got != want {
		t.Fatalf("unexpected payload\nwant: %s\n got: %s", want, got)
	}
}

func TestCRCIsAlwaysFourCharacters(t *testing.T) {
	got, err := generateVietQR(vietQRRequest{
		BankBIN:       "970415",
		AccountNumber: "0011001932418",
		Amount:        120000,
		Purpose:       "TEST5",
		Dynamic:       true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(got, "63040FF5") {
		t.Fatalf("expected zero-padded CRC suffix, got %s", got[len(got)-8:])
	}
}

func TestBuildQRItem(t *testing.T) {
	item := buildQRItem(paymentRow{
		StudentName: "Nguyen An",
		BankBIN:     "970415",
		BankAccount: "0011001932418",
		Amount:      120000,
		Note:        "HP Nguyen An",
	}, 256)

	if len(item.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", item.Errors)
	}
	if item.VietQR == "" {
		t.Fatal("expected vietqr payload")
	}
	if !strings.HasPrefix(item.QRData, "data:image/png;base64,") {
		t.Fatal("expected png data url")
	}
}

func TestServerAddrDefaultsAwayFromCommonPort(t *testing.T) {
	t.Setenv("PORT", "")
	if got := serverAddr(); got != ":18080" {
		t.Fatalf("expected default server address :18080, got %q", got)
	}
}

func TestPaymentItemsDriveQRAmount(t *testing.T) {
	item := buildQRItem(paymentRow{
		StudentName: "Nguyen An",
		BankBIN:     "970415",
		BankAccount: "0011001932418",
		Amount:      1,
		PaymentItems: []paymentItem{
			{Label: "Tiền học phí Tháng 04", LabelEN: "Tuition fees for April", Amount: 3950000},
			{Label: "Phí xe đưa rước Tháng 04", LabelEN: "Shuttle fees for April", Amount: 3030000},
			{Label: "Tiền học phí Tháng 05", LabelEN: "Tuition fees for May", Amount: 3950000},
		},
		BillNumber: "SUN001",
		Note:       "Hoc phi thang 04+05",
	}, 256)
	if len(item.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", item.Errors)
	}
	if item.Amount != 10930000 {
		t.Fatalf("expected item total to override amount, got %d", item.Amount)
	}
	if !strings.Contains(item.VietQR, "540810930000") {
		t.Fatalf("expected QR amount 10930000 in payload, got %s", item.VietQR)
	}
}

func TestImportCSV(t *testing.T) {
	body := strings.NewReader("student_name,parent_name,class_name,bank_bin,bank_account,email,amount,bill_number,note,tuition_april,shuttle_april\nNguyen An,Nguyen Van Binh,3.02,970415,0011001932418,binh@example.com,120000,SUN001,HP Nguyen An,3950000,3030000\n")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/import/csv", body)
	req.Header.Set("Content-Type", "text/csv")
	rec := httptest.NewRecorder()

	handleImportCSV(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Nguyen An") {
		t.Fatalf("expected imported row, got %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"amount":6980000`) {
		t.Fatalf("expected fee columns to override amount, got %s", rec.Body.String())
	}
}

func TestParsePaymentCSVRowsWithFieldMapping(t *testing.T) {
	body := strings.NewReader("Họ và tên,Phụ huynh,Lớp,BIN ngân hàng,Số tài khoản,Tổng phí,Mã hóa đơn,Nội dung\nNguyen An,Nguyen Van Binh,3.02,970415,0011001932418,120000,SUN001,HP Nguyen An\n")
	rows, err := parseCSVRowsWithMapping(body, map[string]string{
		"Họ và tên":     "student",
		"Phụ huynh":     "parent",
		"Lớp":           "class_name",
		"BIN ngân hàng": "bank_bin",
		"Số tài khoản":  "bank_account",
		"Tổng phí":      "amount",
		"Mã hóa đơn":    "bill_number",
		"Nội dung":      "note",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d", len(rows))
	}
	row := rows[0]
	if row.StudentName != "Nguyen An" || row.ParentName != "Nguyen Van Binh" || row.ClassName != "3.02" {
		t.Fatalf("unexpected mapped identity fields: %+v", row)
	}
	if row.BankBIN != "970415" || row.BankAccount != "0011001932418" || row.Amount != 120000 {
		t.Fatalf("unexpected mapped payment fields: %+v", row)
	}
	if row.BillNumber != "SUN001" || row.Note != "HP Nguyen An" {
		t.Fatalf("unexpected mapped QR metadata: %+v", row)
	}
}

func TestParseXLSXPaymentRowsWithFieldMapping(t *testing.T) {
	data := testXLSX(t, [][]string{
		{"Họ và tên", "Phụ huynh", "Lớp", "BIN ngân hàng", "Số tài khoản", "Tổng phí", "Mã hóa đơn", "Nội dung"},
		{"Nguyen An", "Nguyen Van Binh", "3.02", "970415", "0011001932418", "120000", "SUN001", "HP Nguyen An"},
	})
	table, err := parseImportTableBytes(data, "students.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	rows, err := parsePaymentRows(table, normalizeImportMapping("payments", map[string]string{
		"Họ và tên":     "student",
		"Phụ huynh":     "parent",
		"Lớp":           "class_name",
		"BIN ngân hàng": "bank_bin",
		"Số tài khoản":  "bank_account",
		"Tổng phí":      "amount",
		"Mã hóa đơn":    "bill_number",
		"Nội dung":      "note",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d", len(rows))
	}
	row := rows[0]
	if row.StudentName != "Nguyen An" || row.BankAccount != "0011001932418" || row.Amount != 120000 {
		t.Fatalf("unexpected row parsed from xlsx: %+v", row)
	}
}

func TestRenderPaymentDueEmail(t *testing.T) {
	item := buildQRItem(paymentRow{
		ID:          "row-001",
		StudentName: "Nguyen An",
		ParentName:  "Nguyen Van Binh",
		ClassName:   "3.02",
		BankBIN:     "970415",
		BankAccount: "0011001932418",
		Email:       "parent@example.com",
		Amount:      120000,
		BillNumber:  "SUN001",
		Note:        "HP Nguyen An",
	}, 256)
	email, err := renderPaymentEmail(defaultEmailConfig(), item, "payment_due", "http://localhost:18080", "cid")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(email.HTML, `src="cid:qr-row-001"`) {
		t.Fatalf("expected CID QR image, got %s", email.HTML)
	}
	if email.QRPNGBase64 == "" {
		t.Fatal("expected QR attachment content")
	}
}

func TestRenderPaymentPaidEmailDoesNotAttachQR(t *testing.T) {
	item := buildQRItem(paymentRow{
		ID:          "row-paid-001",
		StudentName: "Nguyen An",
		ParentName:  "Nguyen Van Binh",
		ClassName:   "3.02",
		BankBIN:     "970415",
		BankAccount: "0011001932418",
		Email:       "parent@example.com",
		Amount:      120000,
		BillNumber:  "SUN001",
		Note:        "HP Nguyen An",
	}, 256)
	email, err := renderPaymentEmail(defaultEmailConfig(), item, "payment_paid", "http://localhost:18080", "cid")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(email.HTML, "cid:") {
		t.Fatalf("paid confirmation should not embed QR CID image, got %s", email.HTML)
	}
	if email.QRPNGBase64 != "" || email.QRContentID != "" {
		t.Fatalf("paid confirmation should not attach QR, got content=%q cid=%q", email.QRPNGBase64, email.QRContentID)
	}
	if !strings.Contains(email.Text, "xác nhận thanh toán đã hoàn tất") {
		t.Fatalf("expected paid confirmation text, got %s", email.Text)
	}
}

func TestRenderPaymentDueEmailUsesPaymentItems(t *testing.T) {
	item := buildQRItem(paymentRow{
		ID:          "row-002",
		StudentName: "Nguyen An",
		ParentName:  "Nguyen Van Binh",
		ClassName:   "3.02",
		BankBIN:     "970415",
		BankAccount: "0011001932418",
		Email:       "parent@example.com",
		Amount:      1,
		PaymentItems: []paymentItem{
			{Label: "Tiền học phí Tháng 04", LabelEN: "Tuition fees for April", Amount: 3950000},
			{Label: "Phí xe đưa rước Tháng 04", LabelEN: "Shuttle fees for April", Amount: 3030000},
			{Label: "Tiền học phí Tháng 05", LabelEN: "Tuition fees for May", Amount: 3950000},
		},
		BillNumber: "SUN002",
		Note:       "Hoc phi thang 04+05",
	}, 256)
	email, err := renderPaymentEmail(defaultEmailConfig(), item, "payment_due", "http://localhost:18080", "cid")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Tiền học phí Tháng 04",
		"Shuttle fees for April",
		"Tiền học phí Tháng 05",
		"10.930.000",
	} {
		if !strings.Contains(email.HTML, want) {
			t.Fatalf("expected email HTML to contain %q, got %s", want, email.HTML)
		}
	}
}

func TestEmailCronDailyLimitIsCappedAtGmailLimit(t *testing.T) {
	state := normalizeEmailCronState(emailCronState{
		Enabled:    true,
		DailyLimit: 700,
		SendTime:   "25:99",
		Template:   "",
	})
	if state.DailyLimit != 500 {
		t.Fatalf("expected daily limit capped at 500, got %d", state.DailyLimit)
	}
	if state.SendTime != defaultCronSendTime {
		t.Fatalf("expected invalid send time to reset, got %s", state.SendTime)
	}
	if state.Template != "payment_due" {
		t.Fatalf("expected default template, got %s", state.Template)
	}
}

func TestEmailCronDueOncePerDay(t *testing.T) {
	now := time.Date(2026, 5, 21, 9, 0, 0, 0, time.Local)
	state := normalizeEmailCronState(emailCronState{
		Enabled:    true,
		DailyLimit: 100,
		SendTime:   "08:00",
		Template:   "payment_due",
		Queue: []emailCronJob{
			{ID: "row-001", Status: "queued"},
		},
	})
	if !emailCronDue(state, now) {
		t.Fatal("expected cron to be due after send time")
	}
	state.LastRunDate = localDate(now)
	if emailCronDue(state, now) {
		t.Fatal("expected cron to run once per day")
	}
}

func TestDefaultEmailConfigUsesGmail(t *testing.T) {
	cfg := normalizeEmailConfig(emailConfig{})
	if cfg.Provider != emailProviderGmail {
		t.Fatalf("expected gmail provider by default, got %q", cfg.Provider)
	}
	if err := validateEmailConfigForSend(cfg); err == nil || !strings.Contains(err.Error(), "Gmail address") {
		t.Fatalf("expected missing Gmail address error, got %v", err)
	}

	cfg.GmailAddress = "billing@gmail.com"
	cfg.GmailAppPassword = "abcd efgh ijkl mnop"
	if err := validateEmailConfigForSend(cfg); err != nil {
		t.Fatalf("expected valid Gmail config, got %v", err)
	}
}

func TestBuildGmailMessageIncludesInlineQR(t *testing.T) {
	cfg := defaultEmailConfig()
	cfg.GmailAddress = "billing@gmail.com"
	cfg.GmailAppPassword = "abcd efgh ijkl mnop"
	cfg.From = "ABC SUN <billing@gmail.com>"
	cfg.ReplyTo = "billing@example.edu.vn"
	item := buildQRItem(paymentRow{
		ID:          "row-001",
		StudentName: "Nguyen An",
		ClassName:   "3.02",
		BankBIN:     "970415",
		BankAccount: "0011001932418",
		Email:       "parent@example.com",
		Amount:      120000,
		BillNumber:  "SUN001",
		Note:        "HP Nguyen An",
	}, 256)
	rendered, err := renderPaymentEmail(cfg, item, "payment_due", "http://localhost:18080", "cid")
	if err != nil {
		t.Fatal(err)
	}

	msg, err := buildGmailMessage(cfg, item, rendered, "<test@gmail.com>", time.Date(2026, 5, 21, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	body := string(msg)
	for _, want := range []string{
		"From: ABC SUN <billing@gmail.com>",
		"To: parent@example.com",
		"Reply-To: billing@example.edu.vn",
		"Message-ID: <test@gmail.com>",
		"Content-Type: multipart/related",
		"Content-Type: multipart/alternative",
		"Content-Id: <qr-row-001>",
		"Content-Transfer-Encoding: base64",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected Gmail MIME to contain %q, got %s", want, body)
		}
	}
}

func TestEmailCronRolling24HourQuota(t *testing.T) {
	now := time.Date(2026, 5, 21, 9, 0, 0, 0, time.Local)
	state := normalizeEmailCronState(emailCronState{
		DailyLimit: 500,
		SendHistory: []string{
			now.Add(-23 * time.Hour).Format(time.RFC3339),
			now.Add(-25 * time.Hour).Format(time.RFC3339),
		},
	})
	if got := sentLast24hForState(state, now); got != 1 {
		t.Fatalf("expected one send inside rolling window, got %d", got)
	}
	addEmailSentToState(&state, 2, now)
	if got := sentLast24hForState(state, now); got != 3 {
		t.Fatalf("expected three sends after append, got %d", got)
	}
}

func testXLSX(t *testing.T, rows [][]string) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	writeZipFile(t, writer, "[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
  <Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
</Types>`)
	writeZipFile(t, writer, "_rels/.rels", `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`)
	writeZipFile(t, writer, "xl/workbook.xml", `<?xml version="1.0" encoding="UTF-8"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <sheets><sheet name="Sheet1" sheetId="1" r:id="rId1"/></sheets>
</workbook>`)
	writeZipFile(t, writer, "xl/_rels/workbook.xml.rels", `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
</Relationships>`)
	writeZipFile(t, writer, "xl/worksheets/sheet1.xml", testSheetXML(rows))
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}

func writeZipFile(t *testing.T, writer *zip.Writer, name string, content string) {
	t.Helper()
	file, err := writer.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
}

func testSheetXML(rows [][]string) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	for rowIdx, row := range rows {
		builder.WriteString(`<row r="`)
		builder.WriteString(strconv.Itoa(rowIdx + 1))
		builder.WriteString(`">`)
		for colIdx, value := range row {
			ref := testColumnName(colIdx) + strconv.Itoa(rowIdx+1)
			builder.WriteString(`<c r="`)
			builder.WriteString(ref)
			builder.WriteString(`" t="inlineStr"><is><t>`)
			builder.WriteString(testEscapeXML(value))
			builder.WriteString(`</t></is></c>`)
		}
		builder.WriteString(`</row>`)
	}
	builder.WriteString(`</sheetData></worksheet>`)
	return builder.String()
}

func testColumnName(idx int) string {
	idx++
	name := ""
	for idx > 0 {
		idx--
		name = string(rune('A'+idx%26)) + name
		idx /= 26
	}
	return name
}

func testEscapeXML(value string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(value)
}
