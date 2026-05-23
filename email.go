package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

const (
	emailConfigPath       = "email_config.local.json"
	resendConfigPath      = "resend_config.local.json"
	emailProviderGmail    = "gmail"
	emailProviderResend   = "resend"
	defaultEmailSendDelay = 220 * time.Millisecond
	gmailEmailSendDelay   = 500 * time.Millisecond
)

type emailConfig struct {
	Provider         string `json:"provider"`
	APIKey           string `json:"apiKey,omitempty"`
	GmailAddress     string `json:"gmailAddress,omitempty"`
	GmailAppPassword string `json:"gmailAppPassword,omitempty"`
	From             string `json:"from"`
	ReplyTo          string `json:"replyTo,omitempty"`
	Subject          string `json:"subject"`
	SchoolName       string `json:"schoolName"`
	SchoolNameEN     string `json:"schoolNameEn"`
	PaymentPeriod    string `json:"paymentPeriod"`
	PublicBaseURL    string `json:"publicBaseUrl,omitempty"`
}

type emailConfigResponse struct {
	Provider               string `json:"provider"`
	From                   string `json:"from"`
	ReplyTo                string `json:"replyTo,omitempty"`
	Subject                string `json:"subject"`
	SchoolName             string `json:"schoolName"`
	SchoolNameEN           string `json:"schoolNameEn"`
	PaymentPeriod          string `json:"paymentPeriod"`
	PublicBaseURL          string `json:"publicBaseUrl,omitempty"`
	GmailAddress           string `json:"gmailAddress,omitempty"`
	HasGmailAppPassword    bool   `json:"hasGmailAppPassword"`
	GmailAppPasswordMasked string `json:"gmailAppPasswordMasked,omitempty"`
	HasAPIKey              bool   `json:"hasApiKey"`
	APIKeyMasked           string `json:"apiKeyMasked,omitempty"`
}

type emailBatchRequest struct {
	Rows     []paymentRow `json:"rows"`
	Template string       `json:"template"`
	DryRun   bool         `json:"dryRun"`
	Config   *emailConfig `json:"config,omitempty"`
}

type emailPreviewResponse struct {
	Subject string `json:"subject"`
	HTML    string `json:"html"`
	Text    string `json:"text"`
	To      string `json:"to"`
}

type emailSendResult struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	StudentName string `json:"studentName"`
	Status      string `json:"status"`
	Provider    string `json:"provider,omitempty"`
	ResendID    string `json:"resendId,omitempty"`
	MessageID   string `json:"messageId,omitempty"`
	Error       string `json:"error,omitempty"`
	Transient   bool   `json:"transient,omitempty"`
}

type renderedEmail struct {
	Subject     string
	HTML        string
	Text        string
	QRPNGBase64 string
	QRContentID string
}

type resendEmailRequest struct {
	From        string             `json:"from"`
	To          []string           `json:"to"`
	Subject     string             `json:"subject"`
	HTML        string             `json:"html"`
	Text        string             `json:"text,omitempty"`
	ReplyTo     []string           `json:"reply_to,omitempty"`
	Attachments []resendAttachment `json:"attachments,omitempty"`
	Tags        []resendTag        `json:"tags,omitempty"`
}

type resendAttachment struct {
	Content     string `json:"content"`
	Filename    string `json:"filename"`
	ContentID   string `json:"contentId,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

type resendTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func handleEmailConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cfg, err := loadEmailConfig()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, emailConfigPublic(cfg))
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		var req emailConfig
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		current, _ := loadEmailConfig()
		req = normalizeEmailConfig(req)
		if req.APIKey == "" {
			req.APIKey = current.APIKey
		}
		if req.GmailAppPassword == "" {
			req.GmailAppPassword = current.GmailAppPassword
		}
		if err := saveEmailConfig(req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, emailConfigPublic(req))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleEmailPreview(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4<<20)
	var req emailBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	cfg, err := loadEmailConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.Config != nil {
		cfg = normalizeEmailConfig(*req.Config)
	}
	if len(req.Rows) == 0 {
		http.Error(w, "missing rows", http.StatusBadRequest)
		return
	}
	item := buildQRItem(req.Rows[0], 512)
	if len(item.Errors) > 0 {
		http.Error(w, strings.Join(item.Errors, "; "), http.StatusBadRequest)
		return
	}
	email, err := renderPaymentEmail(cfg, item, req.Template, appBaseURL(r, cfg), "data")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, emailPreviewResponse{
		Subject: email.Subject,
		HTML:    email.HTML,
		Text:    email.Text,
		To:      item.Email,
	})
}

func handleEmailSend(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4<<20)
	var req emailBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if len(req.Rows) > maxRows {
		http.Error(w, fmt.Sprintf("too many rows, max is %d", maxRows), http.StatusBadRequest)
		return
	}

	cfg, err := loadEmailConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !req.DryRun {
		if err := validateEmailConfigForSend(cfg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	baseURL := appBaseURL(r, cfg)
	sentLimit := 0
	if !req.DryRun {
		quota, err := emailSendQuotaStatus(time.Now())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if quota.Remaining <= 0 {
			http.Error(w, fmt.Sprintf("email daily limit reached (%d/%d in 24h)", quota.Sent, quota.Limit), http.StatusTooManyRequests)
			return
		}
		sentLimit = quota.Remaining
	}

	results := sendPaymentEmailRows(r.Context(), cfg, req.Rows, req.Template, baseURL, req.DryRun, sentLimit)
	if !req.DryRun {
		recordEmailCronSent(countSentEmails(results), time.Now())
	}

	writeJSON(w, http.StatusOK, map[string]any{"results": results})
}

func countSentEmails(results []emailSendResult) int {
	count := 0
	for _, result := range results {
		if result.Status == "sent" {
			count++
		}
	}
	return count
}

func sendPaymentEmailRows(ctx context.Context, cfg emailConfig, rows []paymentRow, template, baseURL string, dryRun bool, sentLimit int) []emailSendResult {
	results := make([]emailSendResult, 0, len(rows))
	sent := 0
	for _, row := range rows {
		if sentLimit > 0 && sent >= sentLimit {
			results = append(results, skippedEmailSendResult(row, "daily email limit reached"))
			continue
		}
		result := sendPaymentEmailRow(ctx, cfg, row, template, baseURL, dryRun)
		results = append(results, result)
		if result.Status == "sent" {
			sent++
			sleepEmailSendPace(ctx, cfg)
		}
		if result.Transient {
			break
		}
	}
	return results
}

func sendPaymentEmailRow(ctx context.Context, cfg emailConfig, row paymentRow, template, baseURL string, dryRun bool) emailSendResult {
	item := buildQRItem(row, 512)
	result := emailSendResult{
		ID:          item.ID,
		Email:       item.Email,
		StudentName: item.StudentName,
	}
	if len(item.Errors) > 0 {
		result.Status = "error"
		result.Error = strings.Join(item.Errors, "; ")
		return result
	}
	if item.Email == "" {
		result.Status = "skipped"
		result.Error = "missing recipient email"
		return result
	}
	if !strings.Contains(item.Email, "@") {
		result.Status = "error"
		result.Error = "invalid recipient email"
		return result
	}

	email, err := renderPaymentEmail(cfg, item, template, baseURL, "cid")
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		return result
	}
	if dryRun {
		result.Status = "dry_run"
		return result
	}

	outcome, err := sendRenderedEmail(ctx, cfg, item, email)
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		result.Transient = isTransientEmailError(err)
		return result
	}
	result.Status = "sent"
	result.Provider = outcome.Provider
	result.ResendID = outcome.ResendID
	result.MessageID = outcome.MessageID
	return result
}

func loadEmailConfig() (emailConfig, error) {
	cfg := defaultEmailConfig()
	data, err := os.ReadFile(emailConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			legacyData, legacyErr := os.ReadFile(resendConfigPath)
			if legacyErr != nil {
				if errors.Is(legacyErr, os.ErrNotExist) {
					return cfg, nil
				}
				return cfg, legacyErr
			}
			if err := json.Unmarshal(legacyData, &cfg); err != nil {
				return cfg, err
			}
			cfg.Provider = emailProviderResend
			return normalizeEmailConfig(cfg), nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return normalizeEmailConfig(cfg), nil
}

func saveEmailConfig(cfg emailConfig) error {
	cfg = normalizeEmailConfig(cfg)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(emailConfigPath, data, 0600)
}

func defaultEmailConfig() emailConfig {
	return emailConfig{
		Provider:      emailProviderGmail,
		Subject:       "Thông báo thanh toán học phí - {{student_name}}",
		SchoolName:    "Trường Song Ngữ Á Châu",
		SchoolNameEN:  "Asia Bilingual College",
		PaymentPeriod: "Học phí tháng 04+05.2026 - Xe tháng 04.2026",
	}
}

func normalizeEmailConfig(cfg emailConfig) emailConfig {
	cfg.Provider = strings.ToLower(strings.TrimSpace(cfg.Provider))
	if cfg.Provider == "" {
		if cfg.APIKey != "" && cfg.GmailAddress == "" && cfg.GmailAppPassword == "" {
			cfg.Provider = emailProviderResend
		} else {
			cfg.Provider = emailProviderGmail
		}
	}
	if cfg.Provider != emailProviderGmail && cfg.Provider != emailProviderResend {
		cfg.Provider = emailProviderGmail
	}
	cfg.APIKey = strings.TrimSpace(cfg.APIKey)
	cfg.GmailAddress = strings.TrimSpace(cfg.GmailAddress)
	cfg.GmailAppPassword = strings.TrimSpace(cfg.GmailAppPassword)
	cfg.From = strings.TrimSpace(cfg.From)
	cfg.ReplyTo = strings.TrimSpace(cfg.ReplyTo)
	cfg.Subject = strings.TrimSpace(cfg.Subject)
	cfg.SchoolName = strings.TrimSpace(cfg.SchoolName)
	cfg.SchoolNameEN = strings.TrimSpace(cfg.SchoolNameEN)
	cfg.PaymentPeriod = strings.TrimSpace(cfg.PaymentPeriod)
	cfg.PublicBaseURL = strings.TrimRight(strings.TrimSpace(cfg.PublicBaseURL), "/")
	if cfg.Subject == "" {
		cfg.Subject = defaultEmailConfig().Subject
	}
	if cfg.SchoolName == "" {
		cfg.SchoolName = defaultEmailConfig().SchoolName
	}
	if cfg.SchoolNameEN == "" {
		cfg.SchoolNameEN = defaultEmailConfig().SchoolNameEN
	}
	if cfg.PaymentPeriod == "" {
		cfg.PaymentPeriod = defaultEmailConfig().PaymentPeriod
	}
	return cfg
}

func emailConfigPublic(cfg emailConfig) emailConfigResponse {
	cfg = normalizeEmailConfig(cfg)
	return emailConfigResponse{
		Provider:               cfg.Provider,
		From:                   cfg.From,
		ReplyTo:                cfg.ReplyTo,
		Subject:                cfg.Subject,
		SchoolName:             cfg.SchoolName,
		SchoolNameEN:           cfg.SchoolNameEN,
		PaymentPeriod:          cfg.PaymentPeriod,
		PublicBaseURL:          cfg.PublicBaseURL,
		GmailAddress:           cfg.GmailAddress,
		HasGmailAppPassword:    cfg.GmailAppPassword != "",
		GmailAppPasswordMasked: maskAPIKey(cfg.GmailAppPassword),
		HasAPIKey:              cfg.APIKey != "",
		APIKeyMasked:           maskAPIKey(cfg.APIKey),
	}
}

func skippedEmailSendResult(row paymentRow, message string) emailSendResult {
	item := buildQRItem(row, 512)
	return emailSendResult{
		ID:          item.ID,
		Email:       item.Email,
		StudentName: item.StudentName,
		Status:      "skipped",
		Error:       message,
	}
}

func maskAPIKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 8 {
		return "********"
	}
	return key[:5] + strings.Repeat("*", len(key)-9) + key[len(key)-4:]
}

func renderPaymentEmail(cfg emailConfig, item qrItem, templateName, baseURL, imageMode string) (renderedEmail, error) {
	cfg = normalizeEmailConfig(cfg)
	templateName = strings.TrimSpace(templateName)
	if templateName == "" {
		templateName = "payment_due"
	}
	if imageMode != "cid" {
		imageMode = "data"
	}

	qrPNG, err := qrcode.Encode(item.VietQR, qrcode.Medium, 512)
	if err != nil {
		return renderedEmail{}, err
	}
	qrBase64 := base64.StdEncoding.EncodeToString(qrPNG)
	contentID := "qr-" + safeTagValue(item.ID)
	if contentID == "qr-" {
		contentID = "qr-payment"
	}

	qrSrc := "data:image/png;base64," + qrBase64
	if imageMode == "cid" {
		qrSrc = "cid:" + contentID
	}

	qrLink := buildQRLink(baseURL, item)
	subject := renderSubject(cfg.Subject, item)
	switch templateName {
	case "payment_paid":
		return renderedEmail{
			Subject:     subject,
			HTML:        renderPaidEmailHTML(cfg, item),
			Text:        renderPaidEmailText(cfg, item),
			QRPNGBase64: "",
			QRContentID: "",
		}, nil
	case "payment_due":
		return renderedEmail{
			Subject:     subject,
			HTML:        renderPaymentDueHTML(cfg, item, qrSrc, qrLink),
			Text:        renderPaymentDueText(cfg, item, qrLink),
			QRPNGBase64: qrBase64,
			QRContentID: contentID,
		}, nil
	default:
		return renderedEmail{}, fmt.Errorf("unsupported email template %q", templateName)
	}
}

func renderPaymentDueHTML(cfg emailConfig, item qrItem, qrSrc, qrLink string) string {
	student := html.EscapeString(item.StudentName)
	className := html.EscapeString(firstNonEmpty(item.ClassName, "-"))
	period := html.EscapeString(cfg.PaymentPeriod)
	amount := html.EscapeString(formatVND(item.Amount))
	rowsHTML := renderPaymentItemRows(item)
	qrRowspan := len(paymentItemsForEmail(item)) + 2

	return `<!doctype html>
<html><body style="margin:0;padding:0;background:#ffffff;color:#111;font-family:Arial,Helvetica,sans-serif;font-size:16px;line-height:1.35;">
<div style="max-width:760px;margin:0 auto;padding:24px 18px;">
<p style="margin:0 0 24px;">Kính gửi: Quý phụ huynh học sinh <strong>` + student + `</strong> - Lớp: ` + className + `</p>
<p style="margin:0 0 18px;"><em>Dear parent of student <strong>` + student + `</strong> - Class: ` + className + `</em></p>
<p style="margin:0 0 22px;">Trường xin gửi thông báo đến quý phụ huynh các khoản phí cần thanh toán như sau: <em>` + html.EscapeString(cfg.SchoolNameEN) + ` would like to announce the tuition payment that needs to be complete as follows:</em></p>
<table role="presentation" cellpadding="0" cellspacing="0" style="width:100%;border-collapse:collapse;border:1px solid #777;font-size:16px;">
<tr>
<th style="border:1px solid #777;background:#e8e8fb;padding:8px;text-align:center;width:42%;">Diễn giải<br><em>Explanation</em></th>
<th style="border:1px solid #777;background:#e8e8fb;padding:8px;text-align:center;width:24%;">Số tiền<br><em>Amount</em></th>
<th style="border:1px solid #777;background:#e8e8fb;padding:8px;text-align:center;width:34%;">QR</th>
</tr>
<tr>
<td style="border:1px solid #777;padding:8px;">Nội dung thanh toán<br><em>/Payment remark</em></td>
<td style="border:1px solid #777;padding:8px;text-align:center;font-weight:700;">` + period + `</td>
<td rowspan="` + fmt.Sprintf("%d", qrRowspan) + `" style="border:1px solid #777;padding:12px;text-align:center;vertical-align:middle;">
<img alt="QR thanh toán" src="` + qrSrc + `" style="display:block;width:220px;max-width:100%;height:auto;margin:0 auto 10px;border:0;">
<a href="` + html.EscapeString(qrLink) + `" style="color:#1f66c2;text-decoration:underline;">QR Link</a>
</td>
</tr>
` + rowsHTML + `
<tr>
<td style="border:1px solid #777;padding:8px;font-weight:700;">Tổng cộng/<em>Total</em></td>
<td style="border:1px solid #777;padding:8px;text-align:right;font-weight:700;font-size:20px;">` + amount + `</td>
</tr>
</table>
<p style="border:1px solid #777;border-top:0;margin:0;padding:12px;">Quý phụ huynh có thể quét mã QR này để thực hiện thanh toán, hoặc chuyển khoản vào tài khoản định danh của học sinh/<em>Parents can scan this QR code to make payment, or transfer to the student's account.</em></p>
<p style="margin:22px 0 0;font-weight:700;">` + html.EscapeString(cfg.SchoolName) + `</p>
</div>
</body></html>`
}

func renderPaymentItemRows(item qrItem) string {
	var out strings.Builder
	for _, paymentItem := range paymentItemsForEmail(item) {
		label := html.EscapeString(paymentItem.Label)
		labelEN := html.EscapeString(paymentItem.LabelEN)
		out.WriteString(`<tr>
<td style="border:1px solid #777;padding:8px;">` + label + `<br><em>/` + labelEN + `</em></td>
<td style="border:1px solid #777;padding:8px;text-align:right;">` + html.EscapeString(formatVND(paymentItem.Amount)) + `</td>
</tr>`)
	}
	return out.String()
}

func paymentItemsForEmail(item qrItem) []paymentItem {
	if len(item.PaymentItems) > 0 {
		return item.PaymentItems
	}
	return []paymentItem{
		{Label: "Tổng phí cần thanh toán", LabelEN: "Total fees due", Amount: item.Amount},
	}
}

func renderPaidEmailHTML(cfg emailConfig, item qrItem) string {
	student := html.EscapeString(item.StudentName)
	className := html.EscapeString(firstNonEmpty(item.ClassName, "-"))
	amount := html.EscapeString(formatVND(item.Amount))
	today := time.Now().Format("2006-01-02")
	return `<!doctype html>
<html><body style="margin:0;padding:0;background:#ffffff;color:#111;font-family:Arial,Helvetica,sans-serif;font-size:16px;line-height:1.35;">
<div style="max-width:760px;margin:0 auto;padding:24px 18px;">
<p style="margin:0 0 24px;">Kính gửi: Quý phụ huynh học sinh <strong>` + student + `</strong> - Lớp: ` + className + `</p>
<p style="margin:0 0 18px;"><em>Dear parent of student <strong>` + student + `</strong> - Class: ` + className + `</em></p>
<p style="margin:0 0 18px;">Xin cảm ơn quý phụ huynh đã lựa chọn ` + html.EscapeString(cfg.SchoolName) + ` là nơi phát triển tri thức cho con em mình/<em>Thanks for choosing ` + html.EscapeString(cfg.SchoolNameEN) + ` as a place to broaden your child's knowledge.</em></p>
<p style="margin:0 0 22px;">Trường xin gửi thông báo quý phụ huynh đã hoàn tất thanh toán các khoản phí như sau/<em>` + html.EscapeString(cfg.SchoolNameEN) + ` would like to confirm that the following payment has been completed:</em></p>
<table role="presentation" cellpadding="0" cellspacing="0" style="width:100%;border-collapse:collapse;border:1px solid #777;font-size:16px;">
<tr><td style="border:1px solid #777;padding:8px;">Nội dung thanh toán/<br><em>Payment remark</em></td><td style="border:1px solid #777;padding:8px;">` + html.EscapeString(item.Note) + `</td></tr>
<tr><td style="border:1px solid #777;padding:8px;">Lớp/ <em>Classname</em></td><td style="border:1px solid #777;padding:8px;">` + className + `</td></tr>
<tr><td style="border:1px solid #777;padding:8px;">Tài khoản định danh/<em>VA code</em></td><td style="border:1px solid #777;padding:8px;">` + html.EscapeString(item.BillNumber) + `</td></tr>
<tr><td style="border:1px solid #777;padding:8px;">Tình trạng TT/<em>Status</em></td><td style="border:1px solid #777;padding:8px;font-weight:700;">Đã thanh toán/<em>Paid</em></td></tr>
<tr><td style="border:1px solid #777;padding:8px;">Ngày thanh toán /<em>Payment date</em></td><td style="border:1px solid #777;padding:8px;">` + today + `</td></tr>
<tr><td style="border:1px solid #777;padding:8px;">Số tiền thanh toán /<em>Amount</em></td><td style="border:1px solid #777;padding:8px;">` + amount + `</td></tr>
</table>
<p style="margin:22px 0 0;font-weight:700;">` + html.EscapeString(cfg.SchoolName) + `</p>
</div>
</body></html>`
}

func renderPaymentDueText(cfg emailConfig, item qrItem, qrLink string) string {
	lines := make([]string, 0, len(paymentItemsForEmail(item)))
	for _, paymentItem := range paymentItemsForEmail(item) {
		lines = append(lines, fmt.Sprintf("- %s / %s: %s", paymentItem.Label, paymentItem.LabelEN, formatVND(paymentItem.Amount)))
	}
	return fmt.Sprintf("Kính gửi phụ huynh học sinh %s - Lớp %s\n%s thông báo khoản phí cần thanh toán: %s\n%s\nTổng cộng: %s\nMã hóa đơn: %s\nNội dung: %s\nQR Link: %s\n",
		item.StudentName,
		firstNonEmpty(item.ClassName, "-"),
		cfg.SchoolName,
		cfg.PaymentPeriod,
		strings.Join(lines, "\n"),
		formatVND(item.Amount),
		item.BillNumber,
		item.Note,
		qrLink,
	)
}

func renderPaidEmailText(cfg emailConfig, item qrItem) string {
	return fmt.Sprintf("Kính gửi phụ huynh học sinh %s - Lớp %s\n%s xác nhận thanh toán đã hoàn tất.\nSố tiền: %s\nMã hóa đơn: %s\n",
		item.StudentName,
		firstNonEmpty(item.ClassName, "-"),
		cfg.SchoolName,
		formatVND(item.Amount),
		item.BillNumber,
	)
}

func sendResendEmail(ctx context.Context, cfg emailConfig, item qrItem, email renderedEmail) (string, error) {
	payload := resendEmailRequest{
		From:    cfg.From,
		To:      []string{item.Email},
		Subject: email.Subject,
		HTML:    email.HTML,
		Text:    email.Text,
		Tags: []resendTag{
			{Name: "student", Value: safeTagValue(item.StudentName)},
			{Name: "bill", Value: safeTagValue(item.BillNumber)},
		},
	}
	if cfg.ReplyTo != "" {
		payload.ReplyTo = []string{cfg.ReplyTo}
	}
	if email.QRPNGBase64 != "" {
		payload.Attachments = []resendAttachment{
			{
				Content:     email.QRPNGBase64,
				Filename:    "vietqr-" + safeTagValue(item.BillNumber) + ".png",
				ContentID:   email.QRContentID,
				ContentType: "image/png",
			},
		}
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "abcsun-vietqr-demo/1.0")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("resend %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var out struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	return out.ID, nil
}

func renderSubject(subject string, item qrItem) string {
	replacements := map[string]string{
		"{{student_name}}": item.StudentName,
		"{{class_name}}":   item.ClassName,
		"{{bill_number}}":  item.BillNumber,
		"{{amount}}":       formatVND(item.Amount),
	}
	for key, value := range replacements {
		subject = strings.ReplaceAll(subject, key, value)
	}
	return strings.TrimSpace(subject)
}

func buildQRLink(baseURL string, item qrItem) string {
	values := url.Values{}
	values.Set("bankBin", item.BankBIN)
	values.Set("account", item.BankAccount)
	values.Set("amount", fmt.Sprintf("%d", item.Amount))
	values.Set("billNumber", item.BillNumber)
	values.Set("note", item.Note)
	return strings.TrimRight(baseURL, "/") + "/api/v1/qr.png?" + values.Encode()
}

func appBaseURL(r *http.Request, cfg emailConfig) string {
	if cfg.PublicBaseURL != "" {
		return cfg.PublicBaseURL
	}
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}
	return scheme + "://" + r.Host
}

func formatVND(amount int) string {
	raw := fmt.Sprintf("%d", amount)
	var out []byte
	for i, r := range reverse(raw) {
		if i > 0 && i%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, byte(r))
	}
	return string(reverse(string(out))) + " đ"
}

func reverse(value string) []rune {
	runes := []rune(value)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return runes
}

func safeTagValue(value string) string {
	value = cleanANS(value, 80)
	value = regexp.MustCompile(`[^A-Za-z0-9_-]`).ReplaceAllString(value, "_")
	value = strings.Trim(value, "_")
	if value == "" {
		return "unknown"
	}
	if len(value) > 80 {
		return value[:80]
	}
	return value
}
