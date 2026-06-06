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
		TenantID:         "tenant-id",
		TenantStatus:     "active",
		MembershipStatus: "active",
	}
	if err := validateRefreshTokenRecord(valid, now); err != nil {
		t.Fatalf("expected valid refresh token record, got %v", err)
	}

	cases := map[string]authTokenRecord{
		"missing":              {},
		"used":                 {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: valid.ExpiresAt, SessionExpiresAt: valid.SessionExpiresAt, UserStatus: "active", TenantID: valid.TenantID, TenantStatus: valid.TenantStatus, MembershipStatus: valid.MembershipStatus, UsedAt: sql.NullTime{Time: now, Valid: true}},
		"revoked":              {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: valid.ExpiresAt, SessionExpiresAt: valid.SessionExpiresAt, UserStatus: "active", TenantID: valid.TenantID, TenantStatus: valid.TenantStatus, MembershipStatus: valid.MembershipStatus, RevokedAt: sql.NullTime{Time: now, Valid: true}},
		"expired":              {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: now.Add(-time.Second), SessionExpiresAt: valid.SessionExpiresAt, UserStatus: "active", TenantID: valid.TenantID, TenantStatus: valid.TenantStatus, MembershipStatus: valid.MembershipStatus},
		"session expired":      {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: valid.ExpiresAt, SessionExpiresAt: now.Add(-time.Second), UserStatus: "active", TenantID: valid.TenantID, TenantStatus: valid.TenantStatus, MembershipStatus: valid.MembershipStatus},
		"inactive user":        {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: valid.ExpiresAt, SessionExpiresAt: valid.SessionExpiresAt, UserStatus: "suspended", TenantID: valid.TenantID, TenantStatus: valid.TenantStatus, MembershipStatus: valid.MembershipStatus},
		"missing tenant":       {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: valid.ExpiresAt, SessionExpiresAt: valid.SessionExpiresAt, UserStatus: "active", TenantStatus: valid.TenantStatus, MembershipStatus: valid.MembershipStatus},
		"suspended tenant":     {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: valid.ExpiresAt, SessionExpiresAt: valid.SessionExpiresAt, UserStatus: "active", TenantID: valid.TenantID, TenantStatus: "suspended", MembershipStatus: valid.MembershipStatus},
		"suspended membership": {ID: valid.ID, SessionID: valid.SessionID, UserID: valid.UserID, ExpiresAt: valid.ExpiresAt, SessionExpiresAt: valid.SessionExpiresAt, UserStatus: "active", TenantID: valid.TenantID, TenantStatus: valid.TenantStatus, MembershipStatus: "suspended"},
	}
	for name, record := range cases {
		if err := validateRefreshTokenRecord(record, now); err == nil {
			t.Fatalf("expected %s refresh token to be rejected", name)
		}
	}

	platformOnly := authTokenRecord{
		ID:               "refresh-token-id",
		SessionID:        "session-id",
		UserID:           "user-id",
		ExpiresAt:        now.Add(time.Hour),
		SessionExpiresAt: now.Add(24 * time.Hour),
		UserStatus:       "active",
		IsPlatformAdmin:  true,
	}
	if err := validateRefreshTokenRecord(platformOnly, now); err != nil {
		t.Fatalf("expected platform-only refresh token record to be valid, got %v", err)
	}
}

func TestAuthenticatedUserHasActiveTenantAllowsActiveOrTrialTenants(t *testing.T) {
	active := authenticatedUser{ActiveTenant: authTenantSummary{ID: "tenant-1", Status: "active", MembershipStatus: "active"}}
	if !authenticatedUserHasActiveTenant(active) {
		t.Fatal("expected active tenant membership to be accepted")
	}
	trial := authenticatedUser{ActiveTenant: authTenantSummary{ID: "tenant-1", Status: "trial", MembershipStatus: "active"}}
	if !authenticatedUserHasActiveTenant(trial) {
		t.Fatal("expected trial tenant membership to be accepted")
	}
	suspended := authenticatedUser{ActiveTenant: authTenantSummary{ID: "tenant-1", Status: "suspended", MembershipStatus: "active"}}
	if authenticatedUserHasActiveTenant(suspended) {
		t.Fatal("expected suspended tenant to be rejected")
	}
	removed := authenticatedUser{ActiveTenant: authTenantSummary{ID: "tenant-1", Status: "active", MembershipStatus: "removed"}}
	if authenticatedUserHasActiveTenant(removed) {
		t.Fatal("expected removed membership to be rejected")
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

func TestValidateAuthTenantSignupInputRequiresTenantAndOwnerContact(t *testing.T) {
	valid := authTenantSignupRequest{
		TenantName:        "School B",
		TenantCode:        "SCHOOL_B",
		InitialSchoolName: "School B",
		InitialSchoolCode: "SCHOOL_B",
		OwnerEmail:        "owner@example.edu.vn",
		OwnerDisplayName:  "Owner",
		Password:          "very-secure-password",
	}
	if err := validateAuthTenantSignupInput(valid); err != nil {
		t.Fatalf("expected valid signup payload, got %v", err)
	}

	invalid := valid
	invalid.OwnerEmail = ""
	if err := validateAuthTenantSignupInput(invalid); err == nil {
		t.Fatal("expected missing owner contact to be rejected")
	}
}

func TestBuildAuthOnboardingSummaryReflectsTenantOwnerState(t *testing.T) {
	user := authenticatedUser{
		IsTenantOwner: true,
		ActiveTenant: authTenantSummary{
			PlanCode:    "free_trial",
			PlanName:    "Free Trial",
			TrialEndsAt: "2026-07-01",
		},
	}
	got := buildAuthOnboardingSummary(user, 1)
	if !got.IsNewTenantOwner {
		t.Fatal("expected tenant owner onboarding flag")
	}
	if got.NeedsInitialSchool {
		t.Fatal("expected tenant with one school to not need initial school")
	}
	if got.ActiveTenantPlanCode != "free_trial" || got.ActiveTenantTrialEnds != "2026-07-01" {
		t.Fatalf("unexpected onboarding summary: %+v", got)
	}
}

func TestPreferredAuthSessionTenantIDSupportsPlatformOnly(t *testing.T) {
	user := authenticatedUser{IsPlatformAdmin: true}
	got, err := preferredAuthSessionTenantID(user)
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Fatalf("expected empty tenant id for platform-only session, got %q", got)
	}
}

func TestPreferredAuthSessionTenantIDUsesActiveTenantWhenPresent(t *testing.T) {
	user := authenticatedUser{
		IsPlatformAdmin: true,
		ActiveTenant:    authTenantSummary{ID: "tenant-1", Status: "active", MembershipStatus: "active"},
	}
	got, err := preferredAuthSessionTenantID(user)
	if err != nil {
		t.Fatal(err)
	}
	if got != "tenant-1" {
		t.Fatalf("expected tenant-1, got %q", got)
	}
}
