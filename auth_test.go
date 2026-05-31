package main

import (
	"crypto/tls"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHashAuthTokenDoesNotReturnRawToken(t *testing.T) {
	token := "raw-browser-token"
	hash := hashAuthToken(token)

	if hash == "" || hash == token {
		t.Fatalf("expected non-empty hash distinct from token, got %q", hash)
	}
	if hash != hashAuthToken(token) {
		t.Fatal("expected token hash to be deterministic")
	}
	if len(hash) != 64 {
		t.Fatalf("expected sha256 hex hash length 64, got %d", len(hash))
	}
}

func TestNormalizeAuthIdentifierSupportsEmailOrPhone(t *testing.T) {
	if got := normalizeAuthIdentifier(" USER@Example.COM "); got != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", got)
	}
	if got := normalizeAuthIdentifier(" 090 123-4567 "); got != "0901234567" {
		t.Fatalf("expected normalized phone, got %q", got)
	}
}

func TestValidateRefreshTokenRecordRejectsInvalidStates(t *testing.T) {
	now := time.Date(2026, 5, 31, 10, 0, 0, 0, time.UTC)
	valid := authTokenRecord{
		ID:               "refresh-token-id",
		SessionID:        "session-id",
		UserID:           "user-id",
		ExpiresAt:        now.Add(time.Hour),
		SessionExpiresAt: now.Add(24 * time.Hour),
		UserStatus:       "active",
	}
	if err := validateRefreshTokenRecord(valid, now); err != nil {
		t.Fatalf("expected valid refresh token record, got %v", err)
	}

	cases := map[string]authTokenRecord{
		"missing":         {},
		"used":            {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: valid.ExpiresAt, SessionExpiresAt: valid.SessionExpiresAt, UserStatus: "active", UsedAt: sql.NullTime{Time: now, Valid: true}},
		"revoked":         {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: valid.ExpiresAt, SessionExpiresAt: valid.SessionExpiresAt, UserStatus: "active", RevokedAt: sql.NullTime{Time: now, Valid: true}},
		"expired":         {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: now.Add(-time.Second), SessionExpiresAt: valid.SessionExpiresAt, UserStatus: "active"},
		"session expired": {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: valid.ExpiresAt, SessionExpiresAt: now.Add(-time.Second), UserStatus: "active"},
		"inactive user":   {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: valid.ExpiresAt, SessionExpiresAt: valid.SessionExpiresAt, UserStatus: "suspended"},
	}
	for name, record := range cases {
		if err := validateRefreshTokenRecord(record, now); err == nil {
			t.Fatalf("expected %s refresh token to be rejected", name)
		}
	}
}

func TestAuthCookieIsHttpOnlySameSiteAndSecureOnTLS(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "https://example.test/", nil)
	req.TLS = &tls.ConnectionState{}

	cookie := authCookie(req, authConfig{CookieSecure: false}, authAccessCookieName, "token", "/", time.Now().Add(time.Hour))

	if !cookie.HttpOnly {
		t.Fatal("expected auth cookie to be HttpOnly")
	}
	if !cookie.Secure {
		t.Fatal("expected TLS request to force Secure auth cookie")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSite=Lax, got %v", cookie.SameSite)
	}
	if cookie.Path != "/" {
		t.Fatalf("expected cookie path /, got %q", cookie.Path)
	}
}
