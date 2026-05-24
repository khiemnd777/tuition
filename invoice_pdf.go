package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/skip2/go-qrcode"
)

func renderInvoicePDF(invoice invoiceDocument, qr qrItem) ([]byte, error) {
	if qr.VietQR == "" {
		return nil, errors.New("invoice QR payload is required")
	}

	var content strings.Builder
	addPDFText(&content, 48, 798, 18, "ABC SUN Invoice Receipt")
	addPDFText(&content, 48, 774, 11, "School: ABC SUN")
	addPDFText(&content, 48, 752, 11, "Invoice: "+invoice.InvoiceCode)
	addPDFText(&content, 48, 734, 10, "Student: "+invoice.StudentName+" ("+invoice.StudentCode+")")
	addPDFText(&content, 48, 718, 10, "Class: "+invoice.ClassName+" | Period: "+invoice.PeriodCode)
	addPDFText(&content, 48, 702, 10, "Status: "+strings.ToUpper(invoice.Status)+" | Issued: "+invoice.IssuedAt.Format("2006-01-02 15:04"))
	if invoice.DueDate != "" {
		addPDFText(&content, 48, 686, 10, "Due date: "+invoice.DueDate)
	}

	addPDFLine(&content, 48, 668, 547, 668)
	addPDFText(&content, 48, 650, 12, "Line items")
	y := 630.0
	for _, item := range invoice.Items {
		addPDFText(&content, 60, y, 9, truncatePDFText(firstNonEmpty(item.LabelVI, item.LabelEN), 54))
		addPDFText(&content, 430, y, 9, formatPDFVND(item.Amount))
		y -= 15
		if y < 410 {
			break
		}
	}
	if len(invoice.Adjustments) > 0 && y > 390 {
		y -= 8
		addPDFText(&content, 48, y, 12, "Adjustments")
		y -= 18
		for _, adjustment := range invoice.Adjustments {
			label := fmt.Sprintf("%s %s", adjustment.LabelVI, adjustment.Reason)
			addPDFText(&content, 60, y, 9, truncatePDFText(label, 54))
			addPDFText(&content, 430, y, 9, formatPDFVND(adjustment.Delta))
			y -= 15
			if y < 340 {
				break
			}
		}
	}

	addPDFLine(&content, 48, 322, 547, 322)
	addPDFText(&content, 48, 302, 11, "Base amount: "+formatPDFVND(invoice.BaseAmount))
	addPDFText(&content, 48, 284, 11, "Adjustments: "+formatPDFVND(invoice.AdjustmentAmount))
	addPDFText(&content, 48, 262, 14, "Total due: "+formatPDFVND(invoice.TotalAmount))
	addPDFText(&content, 48, 238, 10, "VietQR BillNumber: "+invoice.QRBillNumber)
	addPDFText(&content, 48, 222, 10, "VietQR Purpose: "+invoice.QRNote)
	addPDFText(&content, 48, 206, 9, "Bank: "+invoice.CollectionBankBIN+" / "+invoice.CollectionBankAccount)

	addPDFText(&content, 380, 650, 12, "Payment QR")
	if err := addPDFQRCode(&content, qr.VietQR, 378, 455, 150); err != nil {
		return nil, err
	}
	addPDFText(&content, 372, 438, 8, "Scan to pay this invoice")

	return buildSimplePDF(content.String()), nil
}

func addPDFText(content *strings.Builder, x float64, y float64, size int, text string) {
	fmt.Fprintf(content, "BT /F1 %d Tf %.2f %.2f Td (%s) Tj ET\n", size, x, y, escapePDFString(pdfSafeText(text)))
}

func addPDFLine(content *strings.Builder, x1 float64, y1 float64, x2 float64, y2 float64) {
	fmt.Fprintf(content, "0.75 w %.2f %.2f m %.2f %.2f l S\n", x1, y1, x2, y2)
}

func addPDFQRCode(content *strings.Builder, payload string, x float64, y float64, size float64) error {
	code, err := qrcode.New(payload, qrcode.Medium)
	if err != nil {
		return err
	}
	bitmap := code.Bitmap()
	if len(bitmap) == 0 {
		return errors.New("empty QR bitmap")
	}
	module := size / float64(len(bitmap))
	content.WriteString("0 0 0 rg\n")
	for rowIndex, row := range bitmap {
		for colIndex, on := range row {
			if !on {
				continue
			}
			px := x + float64(colIndex)*module
			py := y + size - float64(rowIndex+1)*module
			fmt.Fprintf(content, "%.3f %.3f %.3f %.3f re f\n", px, py, module+0.01, module+0.01)
		}
	}
	content.WriteString("0 0 0 RG\n")
	return nil
}

func buildSimplePDF(content string) []byte {
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
		fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(content), content),
	}

	var out bytes.Buffer
	out.WriteString("%PDF-1.4\n")
	offsets := make([]int, 0, len(objects)+1)
	offsets = append(offsets, 0)
	for idx, object := range objects {
		offsets = append(offsets, out.Len())
		fmt.Fprintf(&out, "%d 0 obj\n%s\nendobj\n", idx+1, object)
	}
	xrefOffset := out.Len()
	fmt.Fprintf(&out, "xref\n0 %d\n", len(objects)+1)
	out.WriteString("0000000000 65535 f \n")
	for _, offset := range offsets[1:] {
		fmt.Fprintf(&out, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&out, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, xrefOffset)
	return out.Bytes()
}

func escapePDFString(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, "(", `\(`)
	value = strings.ReplaceAll(value, ")", `\)`)
	return value
}

func pdfSafeText(value string) string {
	value = cleanANS(value, 0)
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}

func truncatePDFText(value string, maxLen int) string {
	value = pdfSafeText(value)
	if maxLen > 0 && len(value) > maxLen {
		return strings.TrimSpace(value[:maxLen])
	}
	return value
}

func formatPDFVND(amount int) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}
	raw := fmt.Sprintf("%d", amount)
	parts := []string{}
	for len(raw) > 3 {
		parts = append([]string{raw[len(raw)-3:]}, parts...)
		raw = raw[:len(raw)-3]
	}
	parts = append([]string{raw}, parts...)
	return sign + strings.Join(parts, ".") + " VND"
}
