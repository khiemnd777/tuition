package main

import (
	"archive/zip"
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
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

type importTable struct {
	Headers []string
	Header  map[string]int
	Records [][]string
}

type importFileRequest struct {
	Table   importTable
	Mapping map[string]string
}

type importFieldOption struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Required bool   `json:"required,omitempty"`
}

type importPreviewRow struct {
	RowNumber int               `json:"rowNumber"`
	Values    map[string]string `json:"values"`
}

type importFieldsResponse struct {
	Target           string              `json:"target"`
	Headers          []string            `json:"headers"`
	Fields           []importFieldOption `json:"fields"`
	SuggestedMapping map[string]string   `json:"suggestedMapping"`
	Preview          []importPreviewRow  `json:"preview"`
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
	registerAPIRoutes(mux)

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

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
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

func handleImportFields(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4<<20)

	target, err := normalizeImportTarget(r.URL.Query().Get("target"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req, err := readImportFileRequest(r, target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, buildImportFieldsResponse(target, req.Table))
}

func handleImportCSV(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4<<20)

	req, err := readImportFileRequest(r, "payments")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rows, err := parsePaymentRows(req.Table, req.Mapping)
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
	return parseCSVRowsWithMapping(input, nil)
}

func parseCSVRowsWithMapping(input io.Reader, mapping map[string]string) ([]paymentRow, error) {
	table, err := readCSVTable(input)
	if err != nil {
		return nil, err
	}
	return parsePaymentRows(table, normalizeImportMapping("payments", mapping))
}

func parsePaymentRows(table importTable, mapping map[string]string) ([]paymentRow, error) {
	rows := make([]paymentRow, 0, len(table.Records))
	for i, record := range table.Records {
		if isBlankRecord(record) {
			continue
		}
		row := paymentRow{
			ID:           fmt.Sprintf("csv-%03d", i+1),
			StudentName:  importFieldValue(record, table, mapping, "student", csvAliases("student")),
			ParentName:   importFieldValue(record, table, mapping, "parent", csvAliases("parent")),
			ClassName:    importFieldValue(record, table, mapping, "class_name", csvAliases("class_name")),
			BankBIN:      importFieldValue(record, table, mapping, "bank_bin", csvAliases("bank_bin")),
			BankAccount:  importFieldValue(record, table, mapping, "bank_account", csvAliases("bank_account")),
			Email:        importFieldValue(record, table, mapping, "email", csvAliases("email")),
			Amount:       parseAmount(importFieldValue(record, table, mapping, "amount", csvAliases("amount"))),
			PaymentItems: parsePaymentItemsFromTable(record, table, mapping),
			BillNumber:   importFieldValue(record, table, mapping, "bill_number", csvAliases("bill_number")),
			Note:         importFieldValue(record, table, mapping, "note", csvAliases("note")),
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

func parsePaymentItemsFromTable(record []string, table importTable, mapping map[string]string) []paymentItem {
	if raw := importFieldValue(record, table, mapping, "payment_items", csvAliases("payment_items")); raw != "" {
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
		value := importFieldValue(record, table, mapping, def.Key, csvAliases(def.Key))
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

func readImportFileRequest(r *http.Request, target string) (importFileRequest, error) {
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(4 << 20); err != nil {
			return importFileRequest{}, errors.New("invalid multipart body")
		}
		mapping, err := parseImportMapping(r.FormValue("mapping"), target)
		if err != nil {
			return importFileRequest{}, err
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			return importFileRequest{}, errors.New("missing file field")
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			return importFileRequest{}, errors.New("cannot read file")
		}
		table, err := parseImportTableBytes(data, header.Filename)
		if err != nil {
			return importFileRequest{}, err
		}
		return importFileRequest{Table: table, Mapping: mapping}, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return importFileRequest{}, errors.New("cannot read body")
	}
	table, err := parseImportTableBytes(body, "")
	if err != nil {
		return importFileRequest{}, err
	}
	return importFileRequest{Table: table}, nil
}

func parseImportTableBytes(data []byte, filename string) (importTable, error) {
	if isXLSXFile(data, filename) {
		return readXLSXTable(data)
	}
	return readCSVTable(bytes.NewReader(data))
}

func isXLSXFile(data []byte, filename string) bool {
	if strings.EqualFold(path.Ext(filename), ".xlsx") {
		return true
	}
	return len(data) >= 4 && data[0] == 'P' && data[1] == 'K' && data[2] == 3 && data[3] == 4
}

func readCSVTable(input io.Reader) (importTable, error) {
	reader := csv.NewReader(input)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return importTable{}, fmt.Errorf("cannot parse csv: %w", err)
	}
	return newImportTable(records, "empty csv")
}

func readXLSXTable(data []byte) (importTable, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return importTable{}, fmt.Errorf("cannot parse xlsx: %w", err)
	}
	sharedStrings, err := readXLSXSharedStrings(zipReader)
	if err != nil {
		return importTable{}, err
	}
	sheetPath, err := firstXLSXWorksheetPath(zipReader)
	if err != nil {
		return importTable{}, err
	}
	sheet, err := openXLSXFile(zipReader, sheetPath)
	if err != nil {
		return importTable{}, err
	}
	defer sheet.Close()

	records, err := readXLSXSheetRows(sheet, sharedStrings)
	if err != nil {
		return importTable{}, err
	}
	return newImportTable(records, "empty xlsx")
}

func newImportTable(records [][]string, emptyMessage string) (importTable, error) {
	if len(records) == 0 {
		return importTable{}, errors.New(emptyMessage)
	}
	headers := make([]string, len(records[0]))
	copy(headers, records[0])
	header := map[string]int{}
	for idx, name := range headers {
		key := headerKey(name)
		if key != "" {
			header[key] = idx
		}
	}
	return importTable{
		Headers: headers,
		Header:  header,
		Records: records[1:],
	}, nil
}

func firstXLSXWorksheetPath(zipReader *zip.Reader) (string, error) {
	workbook, err := openXLSXFile(zipReader, "xl/workbook.xml")
	if err != nil {
		if fallback := findXLSXFile(zipReader, "xl/worksheets/sheet1.xml"); fallback != nil {
			return "xl/worksheets/sheet1.xml", nil
		}
		return "", errors.New("xlsx workbook is missing xl/workbook.xml")
	}
	defer workbook.Close()

	decoder := xml.NewDecoder(workbook)
	relationshipID := ""
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("cannot parse xlsx workbook: %w", err)
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "sheet" {
			continue
		}
		for _, attr := range start.Attr {
			if attr.Name.Local == "id" {
				relationshipID = attr.Value
				break
			}
		}
		break
	}
	if relationshipID == "" {
		return "xl/worksheets/sheet1.xml", nil
	}

	rels, err := openXLSXFile(zipReader, "xl/_rels/workbook.xml.rels")
	if err != nil {
		return "", errors.New("xlsx workbook relationships are missing")
	}
	defer rels.Close()

	decoder = xml.NewDecoder(rels)
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("cannot parse xlsx workbook relationships: %w", err)
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "Relationship" {
			continue
		}
		var id, target string
		for _, attr := range start.Attr {
			switch attr.Name.Local {
			case "Id":
				id = attr.Value
			case "Target":
				target = attr.Value
			}
		}
		if id == relationshipID && target != "" {
			return normalizeXLSXPath("xl", target), nil
		}
	}
	return "", errors.New("xlsx first worksheet relationship not found")
}

func normalizeXLSXPath(base string, target string) string {
	target = strings.TrimSpace(strings.ReplaceAll(target, "\\", "/"))
	if strings.HasPrefix(target, "/") {
		return strings.TrimPrefix(path.Clean(target), "/")
	}
	return path.Clean(path.Join(base, target))
}

func readXLSXSharedStrings(zipReader *zip.Reader) ([]string, error) {
	file, err := openXLSXFile(zipReader, "xl/sharedStrings.xml")
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	values := []string{}
	var builder strings.Builder
	inString := false
	inText := false
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot parse xlsx shared strings: %w", err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "si" {
				builder.Reset()
				inString = true
			}
			if inString && t.Name.Local == "t" {
				inText = true
			}
		case xml.CharData:
			if inText {
				builder.Write([]byte(t))
			}
		case xml.EndElement:
			if t.Name.Local == "t" {
				inText = false
			}
			if t.Name.Local == "si" {
				values = append(values, builder.String())
				inString = false
			}
		}
	}
	return values, nil
}

func readXLSXSheetRows(input io.Reader, sharedStrings []string) ([][]string, error) {
	decoder := xml.NewDecoder(input)
	records := [][]string{}
	row := []string{}
	inRow := false
	inCell := false
	inValue := false
	inInlineText := false
	cellType := ""
	cellColumn := 0
	nextColumn := 0
	var valueBuilder strings.Builder
	var inlineBuilder strings.Builder

	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot parse xlsx worksheet: %w", err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "row":
				inRow = true
				row = []string{}
				nextColumn = 0
			case "c":
				if !inRow {
					continue
				}
				inCell = true
				cellType = ""
				cellColumn = nextColumn
				valueBuilder.Reset()
				inlineBuilder.Reset()
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "r":
						if idx := xlsxColumnIndex(attr.Value); idx >= 0 {
							cellColumn = idx
						}
					case "t":
						cellType = attr.Value
					}
				}
			case "v":
				if inCell {
					inValue = true
				}
			case "t":
				if inCell {
					inInlineText = true
				}
			}
		case xml.CharData:
			if inValue {
				valueBuilder.Write([]byte(t))
			}
			if inInlineText {
				inlineBuilder.Write([]byte(t))
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "v":
				inValue = false
			case "t":
				inInlineText = false
			case "c":
				if inCell {
					value := resolveXLSXCellValue(strings.TrimSpace(valueBuilder.String()), inlineBuilder.String(), cellType, sharedStrings)
					row = setCellValue(row, cellColumn, value)
					if cellColumn >= nextColumn {
						nextColumn = cellColumn + 1
					}
				}
				inCell = false
			case "row":
				records = append(records, trimTrailingEmptyCells(row))
				inRow = false
			}
		}
	}
	return records, nil
}

func resolveXLSXCellValue(raw string, inline string, cellType string, sharedStrings []string) string {
	switch cellType {
	case "s":
		idx, err := strconv.Atoi(raw)
		if err == nil && idx >= 0 && idx < len(sharedStrings) {
			return sharedStrings[idx]
		}
		return raw
	case "inlineStr":
		return inline
	case "b":
		if raw == "1" {
			return "true"
		}
		if raw == "0" {
			return "false"
		}
	}
	if inline != "" {
		return inline
	}
	return raw
}

func setCellValue(row []string, idx int, value string) []string {
	for len(row) <= idx {
		row = append(row, "")
	}
	row[idx] = value
	return row
}

func trimTrailingEmptyCells(row []string) []string {
	end := len(row)
	for end > 0 && strings.TrimSpace(row[end-1]) == "" {
		end--
	}
	return row[:end]
}

func xlsxColumnIndex(ref string) int {
	column := 0
	found := false
	for _, r := range ref {
		switch {
		case r >= 'A' && r <= 'Z':
			column = column*26 + int(r-'A'+1)
			found = true
		case r >= 'a' && r <= 'z':
			column = column*26 + int(r-'a'+1)
			found = true
		default:
			if found {
				return column - 1
			}
			return -1
		}
	}
	if !found {
		return -1
	}
	return column - 1
}

func openXLSXFile(zipReader *zip.Reader, name string) (io.ReadCloser, error) {
	file := findXLSXFile(zipReader, name)
	if file == nil {
		return nil, fs.ErrNotExist
	}
	return file.Open()
}

func findXLSXFile(zipReader *zip.Reader, name string) *zip.File {
	for _, file := range zipReader.File {
		if file.Name == name {
			return file
		}
	}
	return nil
}

func importFieldValue(record []string, table importTable, mapping map[string]string, canonical string, aliases []string) string {
	if len(mapping) > 0 {
		for idx, source := range table.Headers {
			if mapping[headerKey(source)] == canonical && idx < len(record) {
				return strings.TrimSpace(record[idx])
			}
		}
	}
	for _, key := range aliases {
		idx, ok := table.Header[key]
		if !ok || idx >= len(record) {
			continue
		}
		if len(mapping) > 0 {
			if _, explicitlyMapped := mapping[headerKey(table.Headers[idx])]; explicitlyMapped {
				continue
			}
		}
		return strings.TrimSpace(record[idx])
	}
	return ""
}

func parseImportMapping(raw string, target string) (map[string]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var mapping map[string]string
	if err := json.Unmarshal([]byte(raw), &mapping); err != nil {
		return nil, errors.New("invalid field mapping")
	}
	return normalizeImportMapping(target, mapping), nil
}

func normalizeImportMapping(target string, mapping map[string]string) map[string]string {
	if len(mapping) == 0 {
		return nil
	}
	normalized := map[string]string{}
	for source, field := range mapping {
		sourceKey := headerKey(source)
		if sourceKey == "" {
			continue
		}
		normalized[sourceKey] = normalizeImportTargetField(target, field)
	}
	return normalized
}

func normalizeImportTarget(raw string) (string, error) {
	switch headerKey(raw) {
	case "", "payment", "payments", "qr", "vietqr", "bang_thanh_toan", "bang_hoc_phi":
		return "payments", nil
	case "master", "master_data", "students", "student", "hoc_sinh", "danh_sach_hoc_sinh":
		return "master_data", nil
	default:
		return "", errors.New("unsupported import target")
	}
}

func normalizeImportTargetField(target string, raw string) string {
	key := headerKey(raw)
	switch key {
	case "", "ignore", "skip", "bo_qua":
		return ""
	}
	for _, field := range importFieldsForTarget(target) {
		if importKeyMatches(key, field.Key) {
			return field.Key
		}
		for _, alias := range importAliasesForTarget(target, field.Key) {
			if importKeyMatches(key, alias) {
				return field.Key
			}
		}
	}
	return key
}

func importKeyMatches(got string, want string) bool {
	want = headerKey(want)
	return got == want || strings.ReplaceAll(got, "_", "") == strings.ReplaceAll(want, "_", "")
}

func buildImportFieldsResponse(target string, table importTable) importFieldsResponse {
	fields := importFieldsForTarget(target)
	suggested := map[string]string{}
	for _, header := range table.Headers {
		suggested[header] = suggestImportField(target, header)
	}
	return importFieldsResponse{
		Target:           target,
		Headers:          table.Headers,
		Fields:           fields,
		SuggestedMapping: suggested,
		Preview:          buildImportPreview(table, 5),
	}
}

func importFieldsForTarget(target string) []importFieldOption {
	if target == "master_data" {
		return []importFieldOption{
			{Key: "student_code", Label: "Mã học sinh", Required: true},
			{Key: "student", Label: "Họ, tên", Required: true},
			{Key: "school", Label: "Trường/cơ sở"},
			{Key: "school_year", Label: "Năm học", Required: true},
			{Key: "grade", Label: "Khối"},
			{Key: "class_name", Label: "Lớp", Required: true},
			{Key: "parent", Label: "Tên ba mẹ"},
			{Key: "parent_email", Label: "Email phụ huynh"},
			{Key: "parent_primary", Label: "Phụ huynh chính"},
			{Key: "parent_active", Label: "Đang active"},
			{Key: "receives_billing_email", Label: "Nhận email học phí"},
		}
	}
	return []importFieldOption{
		{Key: "student", Label: "Họ, tên", Required: true},
		{Key: "parent", Label: "Tên ba mẹ"},
		{Key: "class_name", Label: "Lớp"},
		{Key: "bank_bin", Label: "BIN ngân hàng", Required: true},
		{Key: "bank_account", Label: "Số tài khoản", Required: true},
		{Key: "email", Label: "Email phụ huynh"},
		{Key: "amount", Label: "Tổng phí"},
		{Key: "payment_items", Label: "Khoản phí JSON"},
		{Key: "tuition_april", Label: "Học phí tháng 04"},
		{Key: "shuttle_april", Label: "Phí xe tháng 04"},
		{Key: "tuition_may", Label: "Học phí tháng 05"},
		{Key: "health_insurance", Label: "Bảo hiểm y tế"},
		{Key: "uniform_fee", Label: "Đồng phục"},
		{Key: "international_material", Label: "Sách CTQT"},
		{Key: "previous_fees", Label: "Phí tháng trước"},
		{Key: "bill_number", Label: "Mã hóa đơn"},
		{Key: "note", Label: "Nội dung chuyển khoản"},
	}
}

func suggestImportField(target string, header string) string {
	key := headerKey(header)
	if key == "" {
		return ""
	}
	for _, field := range importFieldsForTarget(target) {
		if importKeyMatches(key, field.Key) {
			return field.Key
		}
		for _, alias := range importAliasesForTarget(target, field.Key) {
			if importKeyMatches(key, alias) {
				return field.Key
			}
		}
	}
	return ""
}

func importAliasesForTarget(target string, field string) []string {
	if target == "master_data" {
		return masterDataCSVAliases(field)
	}
	return csvAliases(field)
}

func buildImportPreview(table importTable, limit int) []importPreviewRow {
	if limit > len(table.Records) {
		limit = len(table.Records)
	}
	preview := make([]importPreviewRow, 0, limit)
	for idx := 0; idx < limit; idx++ {
		record := table.Records[idx]
		values := map[string]string{}
		for column, header := range table.Headers {
			if column < len(record) {
				values[header] = strings.TrimSpace(record[column])
			} else {
				values[header] = ""
			}
		}
		preview = append(preview, importPreviewRow{RowNumber: idx + 2, Values: values})
	}
	return preview
}

func csvAliases(canonical string) []string {
	switch canonical {
	case "student":
		return []string{"student_name", "student", "ten_hoc_sinh", "hoc_sinh", "ho_va_ten", "ho_ten", "ten_hs"}
	case "parent":
		return []string{"parent_name", "parent", "ten_phu_huynh", "phu_huynh", "ten_ba_me", "ba_me", "ten_bo_me"}
	case "class_name":
		return []string{"class_name", "class", "lop", "ten_lop", "lop_hoc"}
	case "bank_bin":
		return []string{"bank_bin", "bin", "ma_bin", "bank_code", "bin_ngan_hang", "ma_ngan_hang"}
	case "bank_account":
		return []string{"bank_account", "account", "account_number", "tai_khoan_ngan_hang", "stk", "so_tai_khoan"}
	case "email":
		return []string{"email", "mail", "email_phu_huynh"}
	case "amount":
		return []string{"amount", "so_tien", "hoc_phi", "tong_phi", "tong_hoc_phi", "so_tien_thanh_toan"}
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
		return []string{"bill_number", "bill", "invoice", "invoice_number", "ma_hoa_don", "ma_tham_chieu", "ma_hd"}
	case "note":
		return []string{"note", "noi_dung", "memo", "noi_dung_chuyen_khoan", "dien_giai"}
	default:
		return []string{canonical}
	}
}

func headerKey(value string) string {
	value = strings.TrimPrefix(strings.TrimSpace(value), "\ufeff")
	value = vietqr.Ascii(strings.ToLower(value))
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
