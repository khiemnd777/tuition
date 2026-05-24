package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/subiz/vietqr"
)

//go:embed web/*
var webFiles embed.FS

const (
	maxRows         = 500
	defaultHTTPPort = "18080"
)

type paymentRow struct {
	ID           string        `json:"id"`
	StudentName  string        `json:"studentName"`
	ParentName   string        `json:"parentName"`
	ClassName    string        `json:"className"`
	BankBIN      string        `json:"bankBin"`
	BankAccount  string        `json:"bankAccount"`
	Email        string        `json:"email"`
	Amount       int           `json:"amount"`
	PaymentItems []paymentItem `json:"paymentItems,omitempty"`
	BillNumber   string        `json:"billNumber"`
	Note         string        `json:"note"`
}

type paymentItem struct {
	Label   string `json:"label"`
	LabelEN string `json:"labelEn"`
	Amount  int    `json:"amount"`
}

type qrItem struct {
	paymentRow
	BankName string   `json:"bankName"`
	VietQR   string   `json:"vietqr"`
	QRData   string   `json:"qrData"`
	Errors   []string `json:"errors,omitempty"`
}

type bankOption struct {
	BIN       string `json:"bin"`
	Code      string `json:"code"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
}

func main() {
	if handled, err := runDBCommand(context.Background(), os.Args[1:], os.Stdout, os.Stderr); handled {
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	runServer()
}

func runServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/banks", method(http.MethodGet, handleBanks))
	mux.HandleFunc("/api/v1/example", method(http.MethodGet, handleExample))
	mux.HandleFunc("/api/v1/import/csv", method(http.MethodPost, handleImportCSV))
	mux.HandleFunc("/api/v1/master-data/options", method(http.MethodGet, handleMasterDataOptions))
	mux.HandleFunc("/api/v1/master-data/students", method(http.MethodGet, handleMasterDataStudents))
	mux.HandleFunc("/api/v1/master-data/import/csv", method(http.MethodPost, handleMasterDataImportCSV))
	mux.HandleFunc("/api/v1/fee-schedules/options", method(http.MethodGet, handleFeeScheduleOptions))
	mux.HandleFunc("/api/v1/fee-schedules", method(http.MethodGet, handleFeeScheduleList))
	mux.HandleFunc("/api/v1/fee-schedules/preview", method(http.MethodPost, handleFeeSchedulePreview))
	mux.HandleFunc("/api/v1/fee-schedules/save", method(http.MethodPost, handleFeeScheduleSave))
	mux.HandleFunc("/api/v1/invoices/options", method(http.MethodGet, handleInvoiceOptions))
	mux.HandleFunc("/api/v1/invoices", method(http.MethodGet, handleInvoiceList))
	mux.HandleFunc("/api/v1/invoices/preview", method(http.MethodPost, handleInvoicePreview))
	mux.HandleFunc("/api/v1/invoices/generate", method(http.MethodPost, handleInvoiceGenerate))
	mux.HandleFunc("/api/v1/invoices/payment", method(http.MethodGet, handleInvoicePayment))
	mux.HandleFunc("/api/v1/invoices/pdf", method(http.MethodGet, handleInvoicePDF))
	mux.HandleFunc("/api/v1/qr.png", method(http.MethodGet, handleQRPNG))
	mux.HandleFunc("/api/v1/vietqr/batch", method(http.MethodPost, handleBatch))
	mux.HandleFunc("/api/v1/email/config", handleEmailConfig)
	mux.HandleFunc("/api/v1/email/preview", method(http.MethodPost, handleEmailPreview))
	mux.HandleFunc("/api/v1/email/send", method(http.MethodPost, handleEmailSend))
	mux.HandleFunc("/api/v1/email/cron", handleEmailCron)
	mux.HandleFunc("/api/v1/email/cron/run", method(http.MethodPost, handleEmailCronRun))

	static, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(static)))

	addr := serverAddr()
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go startEmailCronScheduler(ctx)

	log.Printf("ABC SUN - QR Generating System: %s", localServerURL(addr))
	log.Fatal(server.ListenAndServe())
}

func serverAddr() string {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		return ":" + defaultHTTPPort
	}
	if strings.Contains(port, ":") {
		return port
	}
	return ":" + port
}

func localServerURL(addr string) string {
	if strings.HasPrefix(addr, ":") {
		return "http://localhost" + addr
	}
	return "http://" + addr
}

func method(want string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != want {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

func handleBanks(w http.ResponseWriter, r *http.Request) {
	banks := make([]bankOption, 0, len(vietqr.VNBankM))
	for _, bank := range vietqr.VNBankM {
		banks = append(banks, bankOption{
			BIN:       bank.BIN,
			Code:      bank.Code,
			ShortName: bank.ShortName,
			Name:      bank.Name,
		})
	}
	sort.Slice(banks, func(i, j int) bool {
		return banks[i].ShortName < banks[j].ShortName
	})
	writeJSON(w, http.StatusOK, map[string]any{"banks": banks})
}

func handleExample(w http.ResponseWriter, r *http.Request) {
	row := paymentRow{
		ID:          "demo-vib",
		StudentName: "Nguyễn Duy Khiêm",
		ParentName:  "Nguyễn Duy Khiêm",
		ClassName:   "3.02",
		BankBIN:     "970441",
		BankAccount: "625704060370690",
		Email:       "",
		Amount:      10930000,
		PaymentItems: []paymentItem{
			{Label: "Tiền học phí Tháng 04", LabelEN: "Tuition fees for April", Amount: 3950000},
			{Label: "Phí xe đưa rước Tháng 04", LabelEN: "Shuttle fees for April", Amount: 3030000},
			{Label: "Tiền học phí Tháng 05", LabelEN: "Tuition fees for May", Amount: 3950000},
			{Label: "Bảo hiểm y tế", LabelEN: "Health Insurance", Amount: 0},
			{Label: "Đồng phục", LabelEN: "Uniform fee", Amount: 0},
			{Label: "Sách CTQT", LabelEN: "International material", Amount: 0},
			{Label: "Các khoản phí tháng trước", LabelEN: "Previous month's fees", Amount: 0},
		},
		BillNumber: "TESTVIB01",
		Note:       "Hoc phi thang 04+05",
	}

	item := buildQRItem(row, 360)
	writeJSON(w, http.StatusOK, item)
}

func handleBatch(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4<<20)

	var req struct {
		Rows []paymentRow `json:"rows"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if len(req.Rows) > maxRows {
		http.Error(w, fmt.Sprintf("too many rows, max is %d", maxRows), http.StatusBadRequest)
		return
	}

	items := make([]qrItem, 0, len(req.Rows))
	for idx, row := range req.Rows {
		if strings.TrimSpace(row.ID) == "" {
			row.ID = fmt.Sprintf("row-%03d", idx+1)
		}
		items = append(items, buildQRItem(row, 360))
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func handleQRPNG(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	amount := parseAmount(query.Get("amount"))
	if query.Get("amount") == "" {
		amount = 120000
	}

	row := paymentRow{
		StudentName: "Demo",
		ParentName:  "Demo",
		BankBIN:     firstNonEmpty(query.Get("bankBin"), "970415"),
		BankAccount: firstNonEmpty(query.Get("account"), "0011001932418"),
		Amount:      amount,
		BillNumber:  query.Get("billNumber"),
		Note:        firstNonEmpty(query.Get("note"), "HP Demo"),
	}

	if errs := validateRow(row); len(errs) > 0 {
		http.Error(w, strings.Join(errs, "; "), http.StatusBadRequest)
		return
	}

	row = cleanRow(row)
	payload, err := generateVietQR(vietQRRequest{
		BankBIN:       row.BankBIN,
		AccountNumber: row.BankAccount,
		Amount:        row.Amount,
		BillNumber:    row.BillNumber,
		Purpose:       row.Note,
		Dynamic:       row.Amount > 0,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	png, err := qrcode.Encode(payload, qrcode.Medium, 512)
	if err != nil {
		http.Error(w, "cannot render qr", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(png)
}

func handleImportCSV(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4<<20)

	var reader io.Reader
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(4 << 20); err != nil {
			http.Error(w, "invalid multipart body", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "missing file field", http.StatusBadRequest)
			return
		}
		defer file.Close()
		reader = file
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}
		reader = bytes.NewReader(body)
	}

	rows, err := parseCSVRows(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"rows": rows})
}

func buildQRItem(row paymentRow, size int) qrItem {
	row = cleanRow(row)
	item := qrItem{paymentRow: row}
	if bank, ok := vietqr.VNBankM[row.BankBIN]; ok {
		item.BankName = bank.ShortName
	}

	if errs := validateRow(row); len(errs) > 0 {
		item.Errors = errs
		return item
	}

	payload, err := generateVietQR(vietQRRequest{
		BankBIN:       row.BankBIN,
		AccountNumber: row.BankAccount,
		Amount:        row.Amount,
		BillNumber:    row.BillNumber,
		Purpose:       row.Note,
		Dynamic:       row.Amount > 0,
	})
	if err != nil {
		item.Errors = []string{err.Error()}
		return item
	}
	png, err := qrcode.Encode(payload, qrcode.Medium, size)
	if err != nil {
		item.Errors = []string{"cannot render qr"}
		return item
	}

	item.VietQR = payload
	item.QRData = "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	return item
}

func cleanRow(row paymentRow) paymentRow {
	row.StudentName = strings.TrimSpace(row.StudentName)
	row.ParentName = strings.TrimSpace(row.ParentName)
	row.ClassName = strings.TrimSpace(row.ClassName)
	row.BankBIN = onlyDigits(row.BankBIN)
	row.BankAccount = cleanAccount(row.BankAccount)
	row.Email = strings.TrimSpace(row.Email)
	row.PaymentItems = cleanPaymentItems(row.PaymentItems)
	if total := paymentItemsTotal(row.PaymentItems); len(row.PaymentItems) > 0 {
		row.Amount = total
	}
	row.BillNumber = cleanANS(row.BillNumber, 25)
	row.Note = cleanANS(row.Note, 25)
	if row.Note == "" {
		row.Note = cleanANS("HP "+row.StudentName, 25)
	}
	if row.BillNumber == "" {
		row.BillNumber = defaultBillNumber(row)
	}
	return row
}

func validateRow(row paymentRow) []string {
	var errs []string
	if len(row.BankBIN) != 6 {
		errs = append(errs, "bankBin must be 6 digits")
	} else if _, ok := vietqr.VNBankM[row.BankBIN]; !ok {
		errs = append(errs, "bankBin is not in VietQR bank list")
	}
	if row.BankAccount == "" {
		errs = append(errs, "bankAccount is required")
	}
	if row.Amount < 0 {
		errs = append(errs, "amount must not be negative")
	}
	for idx, item := range row.PaymentItems {
		if item.Amount < 0 {
			errs = append(errs, fmt.Sprintf("paymentItems[%d].amount must not be negative", idx))
		}
	}
	if len(row.BankAccount) > 19 {
		errs = append(errs, "bankAccount max length is 19 characters")
	}
	if row.Email != "" && !strings.Contains(row.Email, "@") {
		errs = append(errs, "email looks invalid")
	}
	return errs
}

func parseCSVRows(input io.Reader) ([]paymentRow, error) {
	reader := csv.NewReader(input)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot parse csv: %w", err)
	}
	if len(records) == 0 {
		return nil, errors.New("empty csv")
	}

	header := map[string]int{}
	for idx, name := range records[0] {
		header[headerKey(name)] = idx
	}

	rows := make([]paymentRow, 0, len(records)-1)
	for i, record := range records[1:] {
		if isBlankRecord(record) {
			continue
		}
		row := paymentRow{
			ID:           fmt.Sprintf("csv-%03d", i+1),
			StudentName:  csvField(record, header, "student"),
			ParentName:   csvField(record, header, "parent"),
			ClassName:    csvField(record, header, "class_name"),
			BankBIN:      csvField(record, header, "bank_bin"),
			BankAccount:  csvField(record, header, "bank_account"),
			Email:        csvField(record, header, "email"),
			Amount:       parseAmount(csvField(record, header, "amount")),
			PaymentItems: parsePaymentItemsFromCSV(record, header),
			BillNumber:   csvField(record, header, "bill_number"),
			Note:         csvField(record, header, "note"),
		}
		rows = append(rows, cleanRow(row))
	}

	if len(rows) > maxRows {
		return nil, fmt.Errorf("too many rows, max is %d", maxRows)
	}
	return rows, nil
}

func cleanPaymentItems(items []paymentItem) []paymentItem {
	cleaned := make([]paymentItem, 0, len(items))
	for _, item := range items {
		item.Label = strings.TrimSpace(item.Label)
		item.LabelEN = strings.TrimSpace(item.LabelEN)
		if item.Label == "" && item.LabelEN == "" && item.Amount == 0 {
			continue
		}
		if item.Label == "" {
			item.Label = item.LabelEN
		}
		if item.LabelEN == "" {
			item.LabelEN = item.Label
		}
		cleaned = append(cleaned, item)
	}
	return cleaned
}

func paymentItemsTotal(items []paymentItem) int {
	total := 0
	for _, item := range items {
		total += item.Amount
	}
	return total
}

func parsePaymentItemsFromCSV(record []string, header map[string]int) []paymentItem {
	if raw := csvField(record, header, "payment_items"); raw != "" {
		var items []paymentItem
		if err := json.Unmarshal([]byte(raw), &items); err == nil {
			return items
		}
	}

	defs := []struct {
		Key     string
		Label   string
		LabelEN string
	}{
		{"tuition_april", "Tiền học phí Tháng 04", "Tuition fees for April"},
		{"shuttle_april", "Phí xe đưa rước Tháng 04", "Shuttle fees for April"},
		{"tuition_may", "Tiền học phí Tháng 05", "Tuition fees for May"},
		{"health_insurance", "Bảo hiểm y tế", "Health Insurance"},
		{"uniform_fee", "Đồng phục", "Uniform fee"},
		{"international_material", "Sách CTQT", "International material"},
		{"previous_fees", "Các khoản phí tháng trước", "Previous month's fees"},
	}

	items := make([]paymentItem, 0, len(defs))
	for _, def := range defs {
		value := csvField(record, header, def.Key)
		if value == "" {
			continue
		}
		items = append(items, paymentItem{
			Label:   def.Label,
			LabelEN: def.LabelEN,
			Amount:  parseAmount(value),
		})
	}
	return items
}

func csvField(record []string, header map[string]int, canonical string) string {
	for _, key := range csvAliases(canonical) {
		if idx, ok := header[key]; ok && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
	}
	return ""
}

func csvAliases(canonical string) []string {
	switch canonical {
	case "student":
		return []string{"student_name", "student", "ten_hoc_sinh", "hoc_sinh"}
	case "parent":
		return []string{"parent_name", "parent", "ten_phu_huynh", "phu_huynh"}
	case "class_name":
		return []string{"class_name", "class", "lop", "ten_lop"}
	case "bank_bin":
		return []string{"bank_bin", "bin", "ma_bin", "bank_code"}
	case "bank_account":
		return []string{"bank_account", "account", "account_number", "tai_khoan_ngan_hang", "stk"}
	case "email":
		return []string{"email", "mail"}
	case "amount":
		return []string{"amount", "so_tien", "hoc_phi"}
	case "payment_items":
		return []string{"payment_items", "items", "fee_items", "khoan_phi"}
	case "tuition_april":
		return []string{"tuition_april", "hoc_phi_thang_04", "hoc_phi_t4"}
	case "shuttle_april":
		return []string{"shuttle_april", "phi_xe_thang_04", "phi_xe_t4"}
	case "tuition_may":
		return []string{"tuition_may", "hoc_phi_thang_05", "hoc_phi_t5"}
	case "health_insurance":
		return []string{"health_insurance", "bao_hiem_y_te"}
	case "uniform_fee":
		return []string{"uniform_fee", "dong_phuc"}
	case "international_material":
		return []string{"international_material", "sach_ctqt", "sach"}
	case "previous_fees":
		return []string{"previous_fees", "cac_khoan_phi_thang_truoc", "phi_thang_truoc"}
	case "bill_number":
		return []string{"bill_number", "bill", "invoice", "invoice_number", "ma_hoa_don", "ma_tham_chieu"}
	case "note":
		return []string{"note", "noi_dung", "memo"}
	default:
		return []string{canonical}
	}
}

func headerKey(value string) string {
	value = vietqr.Ascii(strings.ToLower(strings.TrimSpace(value)))
	value = strings.NewReplacer(" ", "_", "-", "_", ".", "_").Replace(value)
	value = regexp.MustCompile(`_+`).ReplaceAllString(value, "_")
	return strings.Trim(value, "_")
}

func isBlankRecord(record []string) bool {
	for _, value := range record {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}

func onlyDigits(value string) string {
	var builder strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func cleanAccount(value string) string {
	var builder strings.Builder
	for _, r := range strings.TrimSpace(value) {
		switch {
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func parseAmount(value string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	value = strings.NewReplacer(",", "", ".", "", " ", "", "_", "").Replace(value)
	amount, _ := strconv.Atoi(value)
	return amount
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
