package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	adminReportDatasetClasses      = "classes"
	adminReportDatasetInvoices     = "invoices"
	adminReportDatasetTransactions = "transactions"
)

func handleAdminReportsExport(w http.ResponseWriter, r *http.Request) {
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	filters := parseAdminFilters(r)
	dataset := headerKey(firstNonEmpty(r.URL.Query().Get("dataset"), r.URL.Query().Get("type"), adminReportDatasetInvoices))
	filename, data, err := buildAdminReportCSV(r.Context(), db, filters, dataset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func buildAdminReportCSV(ctx context.Context, db *sql.DB, filters adminFilters, dataset string) (string, []byte, error) {
	invoices, err := listAdminInvoiceRows(ctx, db, filters, 5000)
	if err != nil {
		return "", nil, fmt.Errorf("cannot load report invoices")
	}
	var records [][]string
	switch dataset {
	case adminReportDatasetClasses:
		records = adminClassReportCSVRecords(buildAdminClassReportRows(invoices))
	case adminReportDatasetInvoices:
		records = adminInvoiceReportCSVRecords(invoices)
	case adminReportDatasetTransactions:
		transactions, err := listPaymentTransactions(ctx, db, paymentTransactionListFilters{Limit: 5000})
		if err != nil {
			return "", nil, fmt.Errorf("cannot load report transactions")
		}
		records = adminTransactionReportCSVRecords(filterAdminReportTransactions(invoices, transactions, filters))
	default:
		return "", nil, fmt.Errorf("unsupported report dataset %q", dataset)
	}
	data, err := encodeCSVRecords(records)
	if err != nil {
		return "", nil, err
	}
	filename := fmt.Sprintf("abcsun-%s-%s.csv", dataset, time.Now().Format("20060102"))
	return filename, data, nil
}

func encodeCSVRecords(records [][]string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(&buf)
	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}
	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func adminClassReportCSVRecords(rows []adminClassReportRow) [][]string {
	records := [][]string{{
		"school_year",
		"grade",
		"class_name",
		"invoice_count",
		"student_count",
		"total_amount",
		"paid_amount",
		"outstanding_amount",
		"collection_rate",
		"unpaid_count",
		"partial_count",
		"paid_count",
		"overpaid_count",
		"manual_review_count",
	}}
	for _, row := range rows {
		records = append(records, []string{
			row.SchoolYearCode,
			row.Grade,
			row.ClassName,
			strconv.Itoa(row.InvoiceCount),
			strconv.Itoa(row.StudentCount),
			strconv.Itoa(row.TotalAmount),
			strconv.Itoa(row.PaidAmount),
			strconv.Itoa(row.Outstanding),
			fmt.Sprintf("%.4f", row.CollectionRate),
			strconv.Itoa(row.UnpaidCount),
			strconv.Itoa(row.PartialCount),
			strconv.Itoa(row.PaidCount),
			strconv.Itoa(row.OverpaidCount),
			strconv.Itoa(row.ReviewCount),
		})
	}
	return records
}

func adminInvoiceReportCSVRecords(rows []adminInvoiceReportRow) [][]string {
	records := [][]string{{
		"invoice_code",
		"student_code",
		"student_name",
		"class_name",
		"grade",
		"school_year",
		"period_code",
		"month",
		"issued_at",
		"due_date",
		"status",
		"total_amount",
		"paid_amount",
		"outstanding_amount",
		"bank_bin",
		"bank_account",
		"qr_bill_number",
	}}
	for _, invoice := range rows {
		records = append(records, []string{
			invoice.InvoiceCode,
			invoice.StudentCode,
			invoice.StudentName,
			invoice.ClassName,
			invoice.Grade,
			invoice.SchoolYearCode,
			invoice.PeriodCode,
			optionalIntString(invoice.Month),
			invoice.IssuedAt.Format(time.RFC3339),
			invoice.DueDate,
			invoice.Status,
			strconv.Itoa(invoice.TotalAmount),
			strconv.Itoa(invoice.PaidAmount),
			strconv.Itoa(invoice.OutstandingAmount),
			invoice.CollectionBankBIN,
			invoice.CollectionBankAccount,
			invoice.QRBillNumber,
		})
	}
	return records
}

func adminTransactionReportCSVRecords(rows []paymentTransactionSummary) [][]string {
	records := [][]string{{
		"provider",
		"provider_transaction_id",
		"invoice_code",
		"student_code",
		"student_name",
		"direction",
		"amount",
		"currency",
		"transaction_time",
		"account_number",
		"bank_name",
		"reference_code",
		"status",
		"description",
	}}
	for _, transaction := range rows {
		records = append(records, []string{
			transaction.ProviderCode,
			transaction.ProviderTransactionID,
			transaction.InvoiceCode,
			transaction.StudentCode,
			transaction.StudentName,
			transaction.Direction,
			strconv.Itoa(transaction.Amount),
			transaction.Currency,
			transaction.TransactionTime.Format(time.RFC3339),
			transaction.AccountNumber,
			transaction.BankName,
			transaction.ReferenceCode,
			transaction.Status,
			transaction.Description,
		})
	}
	return records
}

func filterAdminReportTransactions(invoices []adminInvoiceReportRow, transactions []paymentTransactionSummary, filters adminFilters) []paymentTransactionSummary {
	if !adminFiltersConstrainInvoices(filters) {
		return transactions
	}
	invoiceIDs := map[string]bool{}
	for _, invoice := range invoices {
		if invoice.ID != "" {
			invoiceIDs[invoice.ID] = true
		}
	}
	filtered := []paymentTransactionSummary{}
	for _, transaction := range transactions {
		if transaction.InvoiceID != "" && invoiceIDs[transaction.InvoiceID] {
			filtered = append(filtered, transaction)
		}
	}
	return filtered
}

func adminFiltersConstrainInvoices(filters adminFilters) bool {
	return strings.TrimSpace(filters.SchoolYearID) != "" ||
		strings.TrimSpace(filters.ClassID) != "" ||
		strings.TrimSpace(filters.Grade) != "" ||
		strings.TrimSpace(filters.PeriodCode) != "" ||
		filters.Month > 0 ||
		strings.TrimSpace(filters.Status) != ""
}

func optionalIntString(value int) string {
	if value == 0 {
		return ""
	}
	return strconv.Itoa(value)
}
