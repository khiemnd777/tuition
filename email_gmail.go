package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"
)

const (
	gmailSMTPAddress = "smtp.gmail.com:587"
	gmailSMTPHost    = "smtp.gmail.com"
)

type emailSendOutcome struct {
	Provider  string
	ResendID  string
	MessageID string
}

func sendRenderedEmail(ctx context.Context, cfg emailConfig, item qrItem, email renderedEmail) (emailSendOutcome, error) {
	cfg = normalizeEmailConfig(cfg)
	switch cfg.Provider {
	case emailProviderResend:
		resendID, err := sendResendEmail(ctx, cfg, item, email)
		return emailSendOutcome{Provider: emailProviderResend, ResendID: resendID}, err
	case emailProviderGmail:
		messageID, err := sendGmailEmail(ctx, cfg, item, email)
		return emailSendOutcome{Provider: emailProviderGmail, MessageID: messageID}, err
	default:
		return emailSendOutcome{}, fmt.Errorf("unsupported email provider %q", cfg.Provider)
	}
}

func validateEmailConfigForSend(cfg emailConfig) error {
	cfg = normalizeEmailConfig(cfg)
	switch cfg.Provider {
	case emailProviderResend:
		if cfg.APIKey == "" {
			return errors.New("missing Resend API key")
		}
		if cfg.From == "" {
			return errors.New("missing sender email")
		}
	case emailProviderGmail:
		if cfg.GmailAddress == "" {
			return errors.New("missing Gmail address")
		}
		if !strings.Contains(cfg.GmailAddress, "@") {
			return errors.New("invalid Gmail address")
		}
		if cfg.GmailAppPassword == "" {
			return errors.New("missing Gmail app password")
		}
	default:
		return fmt.Errorf("unsupported email provider %q", cfg.Provider)
	}
	return nil
}

func sendGmailEmail(ctx context.Context, cfg emailConfig, item qrItem, email renderedEmail) (string, error) {
	if err := validateEmailConfigForSend(cfg); err != nil {
		return "", err
	}
	messageID := buildGmailMessageID(cfg, item)
	data, err := buildGmailMessage(cfg, item, email, messageID, time.Now())
	if err != nil {
		return "", err
	}

	dialer := net.Dialer{Timeout: 20 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", gmailSMTPAddress)
	if err != nil {
		return "", err
	}
	client, err := smtp.NewClient(conn, gmailSMTPHost)
	if err != nil {
		_ = conn.Close()
		return "", err
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); !ok {
		return "", errors.New("gmail smtp does not advertise STARTTLS")
	}
	if err := client.StartTLS(&tls.Config{ServerName: gmailSMTPHost, MinVersion: tls.VersionTLS12}); err != nil {
		return "", err
	}
	auth := smtp.PlainAuth("", cfg.GmailAddress, cfg.GmailAppPassword, gmailSMTPHost)
	if err := client.Auth(auth); err != nil {
		return "", err
	}
	if err := client.Mail(cfg.GmailAddress); err != nil {
		return "", err
	}
	if err := client.Rcpt(item.Email); err != nil {
		return "", err
	}
	wc, err := client.Data()
	if err != nil {
		return "", err
	}
	if _, err := wc.Write(data); err != nil {
		_ = wc.Close()
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	if err := client.Quit(); err != nil {
		return "", err
	}
	return messageID, nil
}

func buildGmailMessage(cfg emailConfig, item qrItem, email renderedEmail, messageID string, now time.Time) ([]byte, error) {
	cfg = normalizeEmailConfig(cfg)
	var buf bytes.Buffer
	writeMailHeader(&buf, "From", formatMailAddress(firstNonEmpty(cfg.From, cfg.GmailAddress)))
	writeMailHeader(&buf, "To", formatMailAddress(item.Email))
	if cfg.ReplyTo != "" {
		writeMailHeader(&buf, "Reply-To", formatMailAddress(cfg.ReplyTo))
	}
	writeMailHeader(&buf, "Subject", mime.QEncoding.Encode("utf-8", email.Subject))
	writeMailHeader(&buf, "Date", now.Format(time.RFC1123Z))
	writeMailHeader(&buf, "Message-ID", messageID)
	writeMailHeader(&buf, "MIME-Version", "1.0")

	related := multipart.NewWriter(&buf)
	writeMailHeader(&buf, "Content-Type", fmt.Sprintf(`multipart/related; boundary="%s"`, related.Boundary()))
	buf.WriteString("\r\n")

	altHeader := textproto.MIMEHeader{}
	altBoundary := multipart.NewWriter(io.Discard).Boundary()
	altHeader.Set("Content-Type", fmt.Sprintf(`multipart/alternative; boundary="%s"`, altBoundary))
	altPart, err := related.CreatePart(altHeader)
	if err != nil {
		return nil, err
	}
	alt := multipart.NewWriter(altPart)
	if err := alt.SetBoundary(altBoundary); err != nil {
		return nil, err
	}
	if err := writeQuotedPrintablePart(alt, "text/plain; charset=UTF-8", email.Text); err != nil {
		return nil, err
	}
	if err := writeQuotedPrintablePart(alt, "text/html; charset=UTF-8", email.HTML); err != nil {
		return nil, err
	}
	if err := alt.Close(); err != nil {
		return nil, err
	}

	if email.QRPNGBase64 != "" {
		imageHeader := textproto.MIMEHeader{}
		filename := "vietqr-" + safeTagValue(item.BillNumber) + ".png"
		imageHeader.Set("Content-Type", fmt.Sprintf(`image/png; name="%s"`, filename))
		imageHeader.Set("Content-Transfer-Encoding", "base64")
		imageHeader.Set("Content-ID", "<"+email.QRContentID+">")
		imageHeader.Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filename))
		imagePart, err := related.CreatePart(imageHeader)
		if err != nil {
			return nil, err
		}
		if err := writeWrappedBase64(imagePart, email.QRPNGBase64); err != nil {
			return nil, err
		}
	}
	if err := related.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeQuotedPrintablePart(writer *multipart.Writer, contentType string, value string) error {
	header := textproto.MIMEHeader{}
	header.Set("Content-Type", contentType)
	header.Set("Content-Transfer-Encoding", "quoted-printable")
	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}
	qp := quotedprintable.NewWriter(part)
	if _, err := qp.Write([]byte(value)); err != nil {
		_ = qp.Close()
		return err
	}
	return qp.Close()
}

func writeWrappedBase64(w io.Writer, value string) error {
	raw := strings.TrimSpace(value)
	for len(raw) > 0 {
		n := 76
		if len(raw) < n {
			n = len(raw)
		}
		if _, err := io.WriteString(w, raw[:n]+"\r\n"); err != nil {
			return err
		}
		raw = raw[n:]
	}
	return nil
}

func buildGmailMessageID(cfg emailConfig, item qrItem) string {
	domain := "gmail.local"
	if idx := strings.LastIndex(cfg.GmailAddress, "@"); idx >= 0 && idx+1 < len(cfg.GmailAddress) {
		domain = cfg.GmailAddress[idx+1:]
	}
	token := safeTagValue(firstNonEmpty(item.ID, item.BillNumber, item.Email))
	return fmt.Sprintf("<%s.%d@%s>", token, time.Now().UnixNano(), domain)
}

func writeMailHeader(buf *bytes.Buffer, key, value string) {
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	buf.WriteString(key)
	buf.WriteString(": ")
	buf.WriteString(value)
	buf.WriteString("\r\n")
}

func formatMailAddress(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	addr, err := mail.ParseAddress(value)
	if err != nil {
		return value
	}
	if addr.Name == "" {
		return addr.Address
	}
	return mime.QEncoding.Encode("utf-8", addr.Name) + " <" + addr.Address + ">"
}

func sleepEmailSendPace(ctx context.Context, cfg emailConfig) {
	delay := emailSendDelay(cfg)
	if delay <= 0 {
		return
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

func emailSendDelay(cfg emailConfig) time.Duration {
	if normalizeEmailConfig(cfg).Provider == emailProviderGmail {
		return gmailEmailSendDelay
	}
	return defaultEmailSendDelay
}

func isTransientEmailError(err error) bool {
	if err == nil {
		return false
	}
	var textErr *textproto.Error
	if errors.As(err, &textErr) {
		return textErr.Code >= 400 && textErr.Code < 500
	}
	message := strings.ToLower(err.Error())
	transientMarkers := []string{
		" 421",
		" 450",
		" 451",
		" 452",
		" 454",
		" 429",
		"rate limit",
		"quota",
		"try again later",
		"too many",
	}
	for _, marker := range transientMarkers {
		if strings.Contains(message, marker) {
			return true
		}
	}
	return false
}
