package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	authAccessCookieName  = "abc_access_token"
	authRefreshCookieName = "abc_refresh_token"

	defaultAccessTokenTTL  = 15 * time.Minute
	defaultRefreshTokenTTL = 7 * 24 * time.Hour
	minAuthPasswordLength  = 12
)

type authContextKey string

const authenticatedUserContextKey authContextKey = "authenticatedUser"

type authConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	CookieSecure    bool
}

type authLoginRequest struct {
	Identifier string `json:"identifier"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type authTenantSummary struct {
	ID                  string `json:"id"`
	Code                string `json:"code"`
	Name                string `json:"name"`
	Status              string `json:"status"`
	MembershipStatus    string `json:"membershipStatus"`
	IsOwner             bool   `json:"isOwner"`
	SubscriptionStatus  string `json:"subscriptionStatus,omitempty"`
	PlanCode            string `json:"planCode,omitempty"`
	PlanName            string `json:"planName,omitempty"`
	TrialEndsAt         string `json:"trialEndsAt,omitempty"`
	CurrentPeriodEndsAt string `json:"currentPeriodEndsAt,omitempty"`
}

type authenticatedUser struct {
	ID                    string                   `json:"id"`
	Email                 string                   `json:"email"`
	Phone                 string                   `json:"phone"`
	DisplayName           string                   `json:"displayName"`
	Status                string                   `json:"status"`
	Tenants               []authTenantSummary      `json:"tenants"`
	ActiveTenant          authTenantSummary        `json:"activeTenant"`
	Roles                 []adminRoleSummary       `json:"roles"`
	ActiveTenantRoles     []adminRoleSummary       `json:"activeTenantRoles"`
	PlatformRoles         []adminRoleSummary       `json:"platformRoles"`
	Permissions           []adminPermissionSummary `json:"permissions"`
	PlatformRoleCodes     []string                 `json:"platformRoleCodes"`
	ActiveTenantRoleCodes []string                 `json:"activeTenantRoleCodes"`
	IsPlatformAdmin       bool                     `json:"isPlatformAdmin"`
	IsTenantOwner         bool                     `json:"isTenantOwner"`
	PermissionSet         map[string]bool          `json:"-"`
}

type authSessionResponse struct {
	User             authenticatedUser     `json:"user"`
	Onboarding       authOnboardingSummary `json:"onboarding"`
	AccessExpiresAt  time.Time             `json:"accessExpiresAt"`
	RefreshExpiresAt time.Time             `json:"refreshExpiresAt"`
}

type authBootstrapRequest struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}

type authBootstrapStatusResponse struct {
	NeedsBootstrap bool `json:"needsBootstrap"`
	SignupEnabled  bool `json:"signupEnabled"`
}

type authTenantSignupRequest struct {
	TenantName        string `json:"tenantName"`
	TenantCode        string `json:"tenantCode"`
	InitialSchoolName string `json:"initialSchoolName"`
	InitialSchoolCode string `json:"initialSchoolCode"`
	OwnerEmail        string `json:"ownerEmail"`
	OwnerPhone        string `json:"ownerPhone"`
	OwnerDisplayName  string `json:"ownerDisplayName"`
	Password          string `json:"password"`
}

type authOnboardingSummary struct {
	IsNewTenantOwner      bool   `json:"isNewTenantOwner"`
	NeedsInitialSchool    bool   `json:"needsInitialSchool"`
	HasMultipleSchools    bool   `json:"hasMultipleSchools"`
	ActiveTenantPlanCode  string `json:"activeTenantPlanCode,omitempty"`
	ActiveTenantPlanName  string `json:"activeTenantPlanName,omitempty"`
	ActiveTenantTrialEnds string `json:"activeTenantTrialEnds,omitempty"`
}

type authTokenRecord struct {
	ID                  string
	SessionID           string
	UserID              string
	ExpiresAt           time.Time
	UsedAt              sql.NullTime
	RevokedAt           sql.NullTime
	SessionExpiresAt    time.Time
	SessionRevokedAt    sql.NullTime
	SessionUserID       string
	SessionIPAddress    string
	SessionUserAgent    string
	UserEmail           string
	UserPhone           string
	UserDisplayName     string
	UserStatus          string
	TenantID            string
	TenantCode          string
	TenantName          string
	TenantStatus        string
	MembershipStatus    string
	TenantIsOwner       bool
	SubscriptionStatus  string
	PlanCode            string
	PlanName            string
	TrialEndsAt         sql.NullTime
	CurrentPeriodEndsAt sql.NullTime
	IsPlatformAdmin     bool
	RefreshExpiresAt    time.Time
	AccessExpiresAt     time.Time
	RefreshTokenID      string
	NewRefreshTokenID   string
	NewAccessTokenID    string
}

type authIssuedTokens struct {
	AccessToken      string
	RefreshToken     string
	AccessTokenHash  string
	RefreshTokenHash string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

func loadAuthConfig() authConfig {
	env, _ := loadAppEnvironment(os.LookupEnv)
	return authConfig{
		AccessTokenTTL:  authEnvDuration("ABC_AUTH_ACCESS_TTL", defaultAccessTokenTTL),
		RefreshTokenTTL: authEnvDuration("ABC_AUTH_REFRESH_TTL", defaultRefreshTokenTTL),
		CookieSecure:    authEnvBool("ABC_AUTH_COOKIE_SECURE", env == appEnvProduction),
	}
}

func authEnvDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func authEnvBool(key string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input authLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	identifier := normalizeAuthIdentifier(firstNonEmpty(input.Identifier, input.Email))
	if identifier == "" || input.Password == "" {
		http.Error(w, "email or phone and password are required", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	if err := ensureBootstrapAdmin(r.Context(), db); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, passwordHash, err := loadAuthUserForLogin(r.Context(), db, identifier)
	if err != nil {
		http.Error(w, "invalid email/phone or password", http.StatusUnauthorized)
		return
	}
	if passwordHash == "" || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)) != nil {
		http.Error(w, "invalid email/phone or password", http.StatusUnauthorized)
		return
	}
	sessionTenantID, err := preferredAuthSessionTenantID(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	now := time.Now().UTC()
	tokens, err := issueAuthSession(r.Context(), db, user.ID, sessionTenantID, now, loadAuthConfig(), auditContextFromRequest(r))
	if err != nil {
		http.Error(w, "cannot create session", http.StatusInternalServerError)
		return
	}
	if err := recordAuthLastLogin(r.Context(), db, user.ID, now); err != nil {
		http.Error(w, "cannot update last login", http.StatusInternalServerError)
		return
	}

	setAuthCookies(w, r, tokens, loadAuthConfig())
	schoolCount := -1
	if sessionTenantID != "" {
		schoolCount = 0
	}
	writeJSON(w, http.StatusOK, authSessionResponse{
		User:             user,
		Onboarding:       buildAuthOnboardingSummary(user, schoolCount),
		AccessExpiresAt:  tokens.AccessExpiresAt,
		RefreshExpiresAt: tokens.RefreshExpiresAt,
	})
}

func handleAuthBootstrap(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleAuthBootstrapStatus(w, r)
	case http.MethodPost:
		handleAuthBootstrapCreate(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAuthBootstrapStatus(w http.ResponseWriter, r *http.Request) {
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	count, err := countAppUsers(r.Context(), db)
	if err != nil {
		http.Error(w, "cannot inspect users", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, authBootstrapStatusResponse{
		NeedsBootstrap: count == 0,
		SignupEnabled:  count > 0,
	})
}

func handleAuthBootstrapCreate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input authBootstrapRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	userInput := adminUserSaveInput{
		Email:       input.Email,
		Phone:       input.Phone,
		DisplayName: input.DisplayName,
		Status:      "active",
		Password:    input.Password,
	}
	if err := validateAdminUserSaveInput(&userInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	userID, err := createInitialAdminUser(r.Context(), db, userInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := authenticatedUser{ID: userID, Email: userInput.Email, Phone: userInput.Phone, DisplayName: userInput.DisplayName, Status: "active"}
	if err := enrichAuthenticatedUser(r.Context(), db, &user); err != nil {
		http.Error(w, "cannot load admin permissions", http.StatusInternalServerError)
		return
	}
	sessionTenantID, err := preferredAuthSessionTenantID(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	now := time.Now().UTC()
	tokens, err := issueAuthSession(r.Context(), db, user.ID, sessionTenantID, now, loadAuthConfig(), auditContextFromRequest(r))
	if err != nil {
		http.Error(w, "cannot create session", http.StatusInternalServerError)
		return
	}
	if err := recordAuthLastLogin(r.Context(), db, user.ID, now); err != nil {
		http.Error(w, "cannot update last login", http.StatusInternalServerError)
		return
	}
	setAuthCookies(w, r, tokens, loadAuthConfig())
	schoolCount := -1
	if sessionTenantID != "" {
		schoolCount = 1
	}
	writeJSON(w, http.StatusCreated, authSessionResponse{
		User:             user,
		Onboarding:       buildAuthOnboardingSummary(user, schoolCount),
		AccessExpiresAt:  tokens.AccessExpiresAt,
		RefreshExpiresAt: tokens.RefreshExpiresAt,
	})
}

func handleAuthTenantSignup(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input authTenantSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizeAuthTenantSignupInput(input)
	if err := validateAuthTenantSignupInput(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	user, schoolCount, err := createTenantOwnerSignup(r.Context(), db, input, auditContextFromRequest(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := enrichAuthenticatedUser(r.Context(), db, &user); err != nil {
		http.Error(w, "cannot load owner session", http.StatusInternalServerError)
		return
	}
	if !authenticatedUserHasActiveTenant(user) {
		http.Error(w, "active tenant membership required", http.StatusForbidden)
		return
	}

	now := time.Now().UTC()
	tokens, err := issueAuthSession(r.Context(), db, user.ID, user.ActiveTenant.ID, now, loadAuthConfig(), auditContextFromRequest(r))
	if err != nil {
		http.Error(w, "cannot create session", http.StatusInternalServerError)
		return
	}
	if err := recordAuthLastLogin(r.Context(), db, user.ID, now); err != nil {
		http.Error(w, "cannot update last login", http.StatusInternalServerError)
		return
	}
	setAuthCookies(w, r, tokens, loadAuthConfig())
	writeJSON(w, http.StatusCreated, authSessionResponse{
		User:             user,
		Onboarding:       buildAuthOnboardingSummary(user, schoolCount),
		AccessExpiresAt:  tokens.AccessExpiresAt,
		RefreshExpiresAt: tokens.RefreshExpiresAt,
	})
}

func handleAuthSession(w http.ResponseWriter, r *http.Request) {
	token, err := readCookieValue(r, authAccessCookieName)
	if err != nil {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	session, err := loadAuthSessionByAccessToken(r.Context(), db, hashAuthToken(token), time.Now().UTC())
	if err != nil {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func handleAuthRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := readCookieValue(r, authRefreshCookieName)
	if err != nil {
		clearAuthCookies(w, r, loadAuthConfig())
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	now := time.Now().UTC()
	session, tokens, err := rotateRefreshToken(r.Context(), db, hashAuthToken(token), now, loadAuthConfig(), auditContextFromRequest(r))
	if err != nil {
		clearAuthCookies(w, r, loadAuthConfig())
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	setAuthCookies(w, r, tokens, loadAuthConfig())
	session.AccessExpiresAt = tokens.AccessExpiresAt
	session.RefreshExpiresAt = tokens.RefreshExpiresAt
	writeJSON(w, http.StatusOK, session)
}

func handleAuthLogout(w http.ResponseWriter, r *http.Request) {
	accessToken, _ := readCookieValue(r, authAccessCookieName)
	refreshToken, _ := readCookieValue(r, authRefreshCookieName)
	if accessToken != "" || refreshToken != "" {
		db, err := openMasterDataDatabase(r.Context())
		if err == nil {
			defer db.Close()
			_ = revokeSessionByTokenHashes(r.Context(), db, hashAuthToken(accessToken), hashAuthToken(refreshToken), "logout")
		}
	}
	clearAuthCookies(w, r, loadAuthConfig())
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func requireAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookieValue(r, authAccessCookieName)
		if err != nil {
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}

		db, err := openMasterDataDatabase(r.Context())
		if err != nil {
			writeMasterDataDBError(w, err)
			return
		}
		defer db.Close()

		session, err := loadAuthSessionByAccessToken(r.Context(), db, hashAuthToken(token), time.Now().UTC())
		if err != nil {
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), authenticatedUserContextKey, session.User)
		next(w, r.WithContext(ctx))
	}
}

func authenticatedUserFromRequest(r *http.Request) (authenticatedUser, bool) {
	user, ok := r.Context().Value(authenticatedUserContextKey).(authenticatedUser)
	return user, ok
}

func activeTenantIDFromRequest(r *http.Request) string {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		return ""
	}
	return strings.TrimSpace(user.ActiveTenant.ID)
}

func authenticatedUserHasActiveTenant(user authenticatedUser) bool {
	if strings.TrimSpace(user.ActiveTenant.ID) == "" {
		return false
	}
	if user.ActiveTenant.Status != "" && !tenantStatusAllowsAuth(user.ActiveTenant.Status) {
		return false
	}
	if user.ActiveTenant.MembershipStatus != "" && user.ActiveTenant.MembershipStatus != "active" {
		return false
	}
	return true
}

func tenantStatusAllowsAuth(status string) bool {
	switch status {
	case "active", "trial":
		return true
	default:
		return false
	}
}

func normalizeAuthEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeAuthIdentifier(value string) string {
	value = strings.TrimSpace(value)
	if strings.Contains(value, "@") {
		return normalizeAuthEmail(value)
	}
	return normalizeAdminPhone(value)
}

func validateAuthPassword(password string) error {
	if len(password) < minAuthPasswordLength {
		return fmt.Errorf("password must be at least %d characters", minAuthPasswordLength)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	if err := validateAuthPassword(password); err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func generateAuthToken() (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func hashAuthToken(token string) string {
	if token == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func issueRawTokens(now time.Time, cfg authConfig) (authIssuedTokens, error) {
	accessToken, err := generateAuthToken()
	if err != nil {
		return authIssuedTokens{}, err
	}
	refreshToken, err := generateAuthToken()
	if err != nil {
		return authIssuedTokens{}, err
	}
	return authIssuedTokens{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessTokenHash:  hashAuthToken(accessToken),
		RefreshTokenHash: hashAuthToken(refreshToken),
		AccessExpiresAt:  now.Add(cfg.AccessTokenTTL),
		RefreshExpiresAt: now.Add(cfg.RefreshTokenTTL),
	}, nil
}

func validateRefreshTokenRecord(record authTokenRecord, now time.Time) error {
	if record.ID == "" || record.SessionID == "" || record.UserID == "" {
		return errors.New("refresh token not found")
	}
	if record.RevokedAt.Valid || record.SessionRevokedAt.Valid {
		return errors.New("refresh token revoked")
	}
	if record.UsedAt.Valid {
		return errors.New("refresh token already used")
	}
	if !record.ExpiresAt.After(now) || !record.SessionExpiresAt.After(now) {
		return errors.New("refresh token expired")
	}
	if record.UserStatus != "active" {
		return errors.New("user is not active")
	}
	if record.TenantID == "" {
		if record.IsPlatformAdmin {
			return nil
		}
		return errors.New("active tenant not found")
	}
	if !tenantStatusAllowsAuth(record.TenantStatus) {
		return errors.New("tenant is not active")
	}
	if record.MembershipStatus != "active" {
		return errors.New("tenant membership is not active")
	}
	return nil
}

func readCookieValue(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return "", http.ErrNoCookie
	}
	return strings.TrimSpace(cookie.Value), nil
}

func setAuthCookies(w http.ResponseWriter, r *http.Request, tokens authIssuedTokens, cfg authConfig) {
	http.SetCookie(w, authCookie(r, cfg, authAccessCookieName, tokens.AccessToken, "/", tokens.AccessExpiresAt))
	http.SetCookie(w, authCookie(r, cfg, authRefreshCookieName, tokens.RefreshToken, "/api/v1/auth", tokens.RefreshExpiresAt))
}

func clearAuthCookies(w http.ResponseWriter, r *http.Request, cfg authConfig) {
	http.SetCookie(w, expiredAuthCookie(r, cfg, authAccessCookieName, "/"))
	http.SetCookie(w, expiredAuthCookie(r, cfg, authRefreshCookieName, "/api/v1/auth"))
}

func authCookie(r *http.Request, cfg authConfig, name string, value string, path string, expires time.Time) *http.Cookie {
	maxAge := int(time.Until(expires).Seconds())
	if maxAge < 1 {
		maxAge = 1
	}
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Expires:  expires,
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   cfg.CookieSecure || r.TLS != nil,
	}
}

func expiredAuthCookie(r *http.Request, cfg authConfig, name string, path string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     path,
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   cfg.CookieSecure || r.TLS != nil,
	}
}

func ensureBootstrapAdmin(ctx context.Context, db *sql.DB) error {
	email := normalizeAuthEmail(os.Getenv("ABC_AUTH_BOOTSTRAP_EMAIL"))
	phone := normalizeAdminPhone(os.Getenv("ABC_AUTH_BOOTSTRAP_PHONE"))
	password := os.Getenv("ABC_AUTH_BOOTSTRAP_PASSWORD")
	displayName := strings.TrimSpace(os.Getenv("ABC_AUTH_BOOTSTRAP_DISPLAY_NAME"))
	if email == "" && phone == "" && password == "" {
		return nil
	}
	if password == "" || (email == "" && phone == "") {
		return fmt.Errorf("ABC_AUTH_BOOTSTRAP_EMAIL or ABC_AUTH_BOOTSTRAP_PHONE must be set with ABC_AUTH_BOOTSTRAP_PASSWORD")
	}
	input := adminUserSaveInput{Email: email, Phone: phone, DisplayName: displayName, Status: "active", Password: password}
	if err := validateAdminUserSaveInput(&input); err != nil {
		return fmt.Errorf("invalid bootstrap admin: %w", err)
	}
	if displayName == "" {
		displayName = input.DisplayName
	}
	passwordHash, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("invalid bootstrap password: %w", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID string
	err = tx.QueryRowContext(ctx, `
SELECT id::text
FROM app_users
WHERE ($1 <> '' AND lower(COALESCE(email, '')) = lower($1))
	OR ($2 <> '' AND phone = $2)
LIMIT 1`, email, phone).Scan(&userID)
	if err == sql.ErrNoRows {
		err = tx.QueryRowContext(ctx, `
INSERT INTO app_users (email, phone, display_name, status, password_hash, password_updated_at)
VALUES ($1, $2, $3, 'active', $4, now())
RETURNING id::text`, email, phone, displayName, passwordHash).Scan(&userID)
	}
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO app_user_roles (user_id, role_id)
SELECT $1::uuid, id
FROM app_roles
WHERE code = 'platform_admin'
ON CONFLICT (user_id, role_id) DO NOTHING`, userID)
	if err != nil {
		return err
	}
	hasRole, err := userHasRole(ctx, tx, userID, "platform_admin")
	if err != nil {
		return err
	}
	if !hasRole {
		return fmt.Errorf("platform_admin role is not available; run migrations first")
	}
	tenantID, err := ensureDefaultTenantMembership(ctx, tx, userID, true)
	if err != nil {
		return err
	}
	if err := ensureTenantUserRole(ctx, tx, tenantID, userID, "tenant_owner", userID); err != nil {
		return err
	}
	return tx.Commit()
}

func createInitialAdminUser(ctx context.Context, db *sql.DB, input adminUserSaveInput) (string, error) {
	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		return "", err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `LOCK TABLE app_users IN SHARE ROW EXCLUSIVE MODE`); err != nil {
		return "", err
	}
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM app_users`).Scan(&count); err != nil {
		return "", err
	}
	if count > 0 {
		return "", fmt.Errorf("admin bootstrap is only available before the first user exists")
	}

	var userID string
	if err := tx.QueryRowContext(ctx, `
INSERT INTO app_users (email, phone, display_name, status, password_hash, password_updated_at)
VALUES ($1, $2, $3, 'active', $4, now())
RETURNING id::text`, input.Email, input.Phone, input.DisplayName, passwordHash).Scan(&userID); err != nil {
		return "", err
	}
	if _, err = tx.ExecContext(ctx, `
INSERT INTO app_user_roles (user_id, role_id)
SELECT $1::uuid, id
FROM app_roles
WHERE code = 'platform_admin'
ON CONFLICT (user_id, role_id) DO NOTHING`, userID); err != nil {
		return "", err
	}
	hasRole, err := userHasRole(ctx, tx, userID, "platform_admin")
	if err != nil {
		return "", err
	}
	if !hasRole {
		return "", fmt.Errorf("platform_admin role is not available; run migrations first")
	}
	tenantID, err := ensureDefaultTenantMembership(ctx, tx, userID, true)
	if err != nil {
		return "", err
	}
	if err := ensureTenantUserRole(ctx, tx, tenantID, userID, "tenant_owner", userID); err != nil {
		return "", err
	}
	return userID, tx.Commit()
}

func normalizeAuthTenantSignupInput(input authTenantSignupRequest) authTenantSignupRequest {
	input.TenantName = strings.TrimSpace(input.TenantName)
	input.TenantCode = normalizeSchoolCode(input.TenantCode)
	input.InitialSchoolName = strings.TrimSpace(input.InitialSchoolName)
	input.InitialSchoolCode = normalizeSchoolCode(input.InitialSchoolCode)
	input.OwnerEmail = strings.ToLower(strings.TrimSpace(input.OwnerEmail))
	input.OwnerPhone = normalizeAdminPhone(input.OwnerPhone)
	input.OwnerDisplayName = strings.TrimSpace(input.OwnerDisplayName)
	input.Password = strings.TrimSpace(input.Password)
	if input.InitialSchoolCode == "" {
		input.InitialSchoolCode = input.TenantCode
	}
	if input.InitialSchoolName == "" {
		input.InitialSchoolName = input.TenantName
	}
	return input
}

func validateAuthTenantSignupInput(input authTenantSignupRequest) error {
	if strings.TrimSpace(input.TenantName) == "" {
		return fmt.Errorf("tenant name is required")
	}
	if !isSafeTenantCode(input.TenantCode) {
		return fmt.Errorf("tenant code may contain only letters, numbers, underscore, or hyphen")
	}
	if !isSafeTenantCode(input.InitialSchoolCode) {
		return fmt.Errorf("initial school code may contain only letters, numbers, underscore, or hyphen")
	}
	if strings.TrimSpace(input.InitialSchoolName) == "" {
		return fmt.Errorf("initial school name is required")
	}
	ownerInput := adminUserSaveInput{
		Email:       input.OwnerEmail,
		Phone:       input.OwnerPhone,
		DisplayName: input.OwnerDisplayName,
		Status:      "active",
		Password:    input.Password,
	}
	return validateAdminUserSaveInput(&ownerInput)
}

func createTenantOwnerSignup(ctx context.Context, db *sql.DB, input authTenantSignupRequest, auditCtx requestAuditContext) (authenticatedUser, int, error) {
	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		return authenticatedUser{}, 0, err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return authenticatedUser{}, 0, err
	}
	defer tx.Rollback()

	var existingUserID string
	err = tx.QueryRowContext(ctx, `
SELECT id::text
FROM app_users
WHERE ($1 <> '' AND lower(COALESCE(email, '')) = lower($1))
	OR ($2 <> '' AND phone = $2)
LIMIT 1`, input.OwnerEmail, input.OwnerPhone).Scan(&existingUserID)
	if err == nil && existingUserID != "" {
		return authenticatedUser{}, 0, fmt.Errorf("owner email or phone already exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return authenticatedUser{}, 0, err
	}

	var user authenticatedUser
	if err := tx.QueryRowContext(ctx, `
INSERT INTO app_users (email, phone, display_name, status, password_hash, password_updated_at)
VALUES ($1, $2, $3, 'active', $4, now())
RETURNING id::text, COALESCE(email, ''), phone, display_name, status`,
		input.OwnerEmail,
		input.OwnerPhone,
		input.OwnerDisplayName,
		passwordHash,
	).Scan(&user.ID, &user.Email, &user.Phone, &user.DisplayName, &user.Status); err != nil {
		return authenticatedUser{}, 0, err
	}

	tenantInput := tenantSaveInput{
		Code:              input.TenantCode,
		Name:              input.TenantName,
		Status:            "trial",
		InitialSchoolCode: input.InitialSchoolCode,
		InitialSchoolName: input.InitialSchoolName,
	}
	tenant, err := createTenantWithInitialSchool(ctx, tx, tenantInput, user, auditCtx)
	if err != nil {
		return authenticatedUser{}, 0, err
	}
	user.ActiveTenant = authTenantSummary{
		ID:                  tenant.ID,
		Code:                tenant.Code,
		Name:                tenant.Name,
		Status:              tenant.Status,
		MembershipStatus:    "active",
		IsOwner:             true,
		SubscriptionStatus:  tenant.SubscriptionStatus,
		PlanCode:            tenant.PlanCode,
		PlanName:            tenant.PlanName,
		TrialEndsAt:         tenant.TrialEndsAt,
		CurrentPeriodEndsAt: tenant.CurrentPeriodEndsAt,
	}
	if err := tx.Commit(); err != nil {
		return authenticatedUser{}, 0, err
	}
	return user, tenant.SchoolCount, nil
}

func userHasRole(ctx context.Context, tx *sql.Tx, userID string, roleCode string) (bool, error) {
	var hasRole bool
	err := tx.QueryRowContext(ctx, `
SELECT EXISTS (
	SELECT 1
	FROM app_user_roles ur
	JOIN app_roles role ON role.id = ur.role_id
	WHERE ur.user_id = $1::uuid
		AND role.code = $2
)`, userID, roleCode).Scan(&hasRole)
	return hasRole, err
}

func ensureDefaultTenantMembership(ctx context.Context, exec masterDataExecutor, userID string, isOwner bool) (string, error) {
	var tenantID string
	err := exec.QueryRowContext(ctx, `
WITH default_tenant AS (
	SELECT id
	FROM tenants
	WHERE code = $2
	LIMIT 1
), upserted AS (
	INSERT INTO tenant_memberships (tenant_id, user_id, status, is_owner)
	SELECT id, $1::uuid, 'active', $3
	FROM default_tenant
	ON CONFLICT (tenant_id, user_id) DO UPDATE
	SET status = 'active',
		is_owner = tenant_memberships.is_owner OR EXCLUDED.is_owner,
		updated_at = now()
	RETURNING tenant_id::text
)
SELECT tenant_id FROM upserted
UNION ALL
SELECT id::text FROM default_tenant
LIMIT 1`, userID, defaultTenantCode, isOwner).Scan(&tenantID)
	return tenantID, err
}

func ensureTenantMembership(ctx context.Context, exec masterDataExecutor, tenantID string, userID string, isOwner bool) error {
	_, err := exec.ExecContext(ctx, `
INSERT INTO tenant_memberships (tenant_id, user_id, status, is_owner)
VALUES ($1::uuid, $2::uuid, 'active', $3)
ON CONFLICT (tenant_id, user_id) DO UPDATE
SET status = 'active',
	is_owner = tenant_memberships.is_owner OR EXCLUDED.is_owner,
	updated_at = now()`, tenantID, userID, isOwner)
	return err
}

func ensureTenantUserRole(ctx context.Context, exec masterDataExecutor, tenantID string, userID string, roleCode string, assignedByUserID string) error {
	_, err := exec.ExecContext(ctx, `
INSERT INTO tenant_user_roles (tenant_id, user_id, role_id, assigned_by_user_id)
SELECT $1::uuid, $2::uuid, role.id, nullif($4, '')::uuid
FROM app_roles role
WHERE role.code = $3
ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING`, tenantID, userID, roleCode, assignedByUserID)
	return err
}

func countAppUsers(ctx context.Context, db *sql.DB) (int, error) {
	var count int
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM app_users`).Scan(&count)
	return count, err
}

func loadAuthUserForLogin(ctx context.Context, db *sql.DB, identifier string) (authenticatedUser, string, error) {
	email := ""
	phone := ""
	if strings.Contains(identifier, "@") {
		email = normalizeAuthEmail(identifier)
	} else {
		phone = normalizeAdminPhone(identifier)
	}
	var user authenticatedUser
	var passwordHash string
	err := db.QueryRowContext(ctx, `
SELECT id::text, COALESCE(email, ''), phone, display_name, status, password_hash
FROM app_users
WHERE (($1 <> '' AND lower(COALESCE(email, '')) = lower($1))
		OR ($2 <> '' AND phone = $2))
	AND status = 'active'
LIMIT 1`, email, phone).Scan(&user.ID, &user.Email, &user.Phone, &user.DisplayName, &user.Status, &passwordHash)
	if err != nil {
		return authenticatedUser{}, "", err
	}
	if err := enrichAuthenticatedUser(ctx, db, &user); err != nil {
		return authenticatedUser{}, "", err
	}
	return user, passwordHash, nil
}

func enrichAuthenticatedUser(ctx context.Context, db *sql.DB, user *authenticatedUser) error {
	tenants, activeTenant, err := loadAuthenticatedUserTenants(ctx, db, user.ID, user.ActiveTenant.ID)
	if err != nil {
		return err
	}
	user.Tenants = tenants
	user.ActiveTenant = activeTenant
	user.Roles = []adminRoleSummary{}
	user.ActiveTenantRoles = []adminRoleSummary{}
	user.PlatformRoles = []adminRoleSummary{}
	user.Permissions = []adminPermissionSummary{}
	user.PlatformRoleCodes = []string{}
	user.ActiveTenantRoleCodes = []string{}
	user.IsPlatformAdmin = false
	user.IsTenantOwner = false
	user.PermissionSet = map[string]bool{}

	platformRoles, platformPermissions, platformRoleCodes, err := loadUserRolesAndPermissions(ctx, db, `
SELECT r.id::text,
	r.code,
	r.name,
	r.description,
	r.is_system,
	COALESCE(p.id::text, ''),
	COALESCE(p.code, ''),
	COALESCE(p.description, '')
FROM app_user_roles ur
JOIN app_roles r ON r.id = ur.role_id
LEFT JOIN app_role_permissions rp ON rp.role_id = r.id
LEFT JOIN app_permissions p ON p.id = rp.permission_id
WHERE ur.user_id = $1::uuid
	AND r.code IN ('platform_admin')
ORDER BY r.code, p.code`, user.ID)
	if err != nil {
		return err
	}
	user.PlatformRoles = platformRoles
	user.PlatformRoleCodes = platformRoleCodes
	user.IsPlatformAdmin = stringSliceContains(platformRoleCodes, "platform_admin")

	tenantRoles := []adminRoleSummary{}
	tenantPermissions := []adminPermissionSummary{}
	tenantRoleCodes := []string{}
	if user.ActiveTenant.ID != "" {
		tenantRoles, tenantPermissions, tenantRoleCodes, err = loadUserRolesAndPermissions(ctx, db, `
SELECT r.id::text,
	r.code,
	r.name,
	r.description,
	r.is_system,
	COALESCE(p.id::text, ''),
	COALESCE(p.code, ''),
	COALESCE(p.description, '')
FROM tenant_user_roles ur
JOIN app_roles r ON r.id = ur.role_id
LEFT JOIN app_role_permissions rp ON rp.role_id = r.id
LEFT JOIN app_permissions p ON p.id = rp.permission_id
WHERE ur.user_id = $1::uuid
	AND ur.tenant_id = $2::uuid
	AND r.code IN ('tenant_owner', 'tenant_admin', 'tenant_staff', 'tenant_accountant')
ORDER BY r.code, p.code`, user.ID, user.ActiveTenant.ID)
		if err != nil {
			return err
		}
	}

	user.Roles = tenantRoles
	user.ActiveTenantRoles = tenantRoles
	user.ActiveTenantRoleCodes = tenantRoleCodes
	user.IsTenantOwner = user.ActiveTenant.IsOwner || stringSliceContains(tenantRoleCodes, "tenant_owner")
	user.Permissions = mergeAdminPermissions(platformPermissions, tenantPermissions)
	for _, permission := range user.Permissions {
		user.PermissionSet[permission.Code] = true
	}
	return nil
}

func loadUserRolesAndPermissions(ctx context.Context, db *sql.DB, query string, args ...any) ([]adminRoleSummary, []adminPermissionSummary, []string, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	roleByID := map[string]*adminRoleSummary{}
	roleOrder := []string{}
	permissionByCode := map[string]adminPermissionSummary{}
	for rows.Next() {
		var roleID string
		var role adminRoleSummary
		var permission adminPermissionSummary
		if err := rows.Scan(
			&roleID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.IsSystem,
			&permission.ID,
			&permission.Code,
			&permission.Description,
		); err != nil {
			return nil, nil, nil, err
		}
		existing := roleByID[roleID]
		if existing == nil {
			role.ID = roleID
			role.Permissions = []adminPermissionSummary{}
			existing = &role
			roleByID[roleID] = existing
			roleOrder = append(roleOrder, roleID)
		}
		if permission.Code != "" {
			existing.Permissions = append(existing.Permissions, permission)
			permissionByCode[permission.Code] = permission
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, err
	}

	roles := make([]adminRoleSummary, 0, len(roleOrder))
	roleCodes := make([]string, 0, len(roleOrder))
	for _, roleID := range roleOrder {
		role := *roleByID[roleID]
		sortAdminPermissions(role.Permissions)
		roles = append(roles, role)
		roleCodes = append(roleCodes, role.Code)
	}
	permissions := make([]adminPermissionSummary, 0, len(permissionByCode))
	for _, permission := range permissionByCode {
		permissions = append(permissions, permission)
	}
	sortAdminPermissions(permissions)
	sort.Strings(roleCodes)
	return roles, permissions, roleCodes, nil
}

func mergeAdminPermissions(groups ...[]adminPermissionSummary) []adminPermissionSummary {
	permissionByCode := map[string]adminPermissionSummary{}
	for _, group := range groups {
		for _, permission := range group {
			if permission.Code == "" {
				continue
			}
			permissionByCode[permission.Code] = permission
		}
	}
	permissions := make([]adminPermissionSummary, 0, len(permissionByCode))
	for _, permission := range permissionByCode {
		permissions = append(permissions, permission)
	}
	sortAdminPermissions(permissions)
	return permissions
}

func stringSliceContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func buildAuthOnboardingSummary(user authenticatedUser, schoolCount int) authOnboardingSummary {
	summary := authOnboardingSummary{
		IsNewTenantOwner:      user.IsTenantOwner,
		NeedsInitialSchool:    schoolCount == 0,
		HasMultipleSchools:    schoolCount > 1,
		ActiveTenantPlanCode:  user.ActiveTenant.PlanCode,
		ActiveTenantPlanName:  user.ActiveTenant.PlanName,
		ActiveTenantTrialEnds: user.ActiveTenant.TrialEndsAt,
	}
	if schoolCount < 0 {
		summary.NeedsInitialSchool = false
		summary.HasMultipleSchools = false
	}
	return summary
}

func loadTenantSchoolCount(ctx context.Context, exec masterDataExecutor, tenantID string) (int, error) {
	if strings.TrimSpace(tenantID) == "" {
		return 0, nil
	}
	var count int
	err := exec.QueryRowContext(ctx, `
SELECT COUNT(*)
FROM schools
WHERE tenant_id = $1::uuid`, tenantID).Scan(&count)
	return count, err
}

func loadAuthenticatedUserTenants(ctx context.Context, db *sql.DB, userID string, preferredTenantID string) ([]authTenantSummary, authTenantSummary, error) {
	rows, err := db.QueryContext(ctx, `
SELECT tenant.id::text,
	COALESCE(tenant.code, ''),
	COALESCE(tenant.name, ''),
	tenant.status,
	membership.status,
	membership.is_owner,
	COALESCE(ts.status, ''),
	COALESCE(plan.code, ''),
	COALESCE(plan.name, ''),
	ts.trial_ends_at,
	ts.current_period_ends_at
FROM tenant_memberships membership
JOIN tenants tenant ON tenant.id = membership.tenant_id
LEFT JOIN tenant_subscriptions ts ON ts.tenant_id = tenant.id
LEFT JOIN subscription_plans plan ON plan.id = ts.plan_id
WHERE membership.user_id = $1::uuid
	AND membership.status = 'active'
	AND tenant.status IN ('active', 'trial')
ORDER BY CASE WHEN tenant.id::text = $2 THEN 0 WHEN tenant.code = $3 THEN 1 ELSE 2 END,
	tenant.code`, userID, strings.TrimSpace(preferredTenantID), defaultTenantCode)
	if err != nil {
		return nil, authTenantSummary{}, err
	}
	defer rows.Close()

	tenants := []authTenantSummary{}
	activeTenant := authTenantSummary{}
	for rows.Next() {
		var item authTenantSummary
		var trialEndsAt sql.NullTime
		var currentPeriodEndsAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Status, &item.MembershipStatus, &item.IsOwner, &item.SubscriptionStatus, &item.PlanCode, &item.PlanName, &trialEndsAt, &currentPeriodEndsAt); err != nil {
			return nil, authTenantSummary{}, err
		}
		item.TrialEndsAt = formatNullDate(trialEndsAt)
		item.CurrentPeriodEndsAt = formatNullDate(currentPeriodEndsAt)
		if activeTenant.ID == "" {
			activeTenant = item
		}
		if preferredTenantID != "" && item.ID == preferredTenantID {
			activeTenant = item
		}
		tenants = append(tenants, item)
	}
	if err := rows.Err(); err != nil {
		return nil, authTenantSummary{}, err
	}
	return tenants, activeTenant, nil
}

func sortAdminPermissions(permissions []adminPermissionSummary) {
	sort.Slice(permissions, func(i, j int) bool {
		return permissions[i].Code < permissions[j].Code
	})
}

func preferredAuthSessionTenantID(user authenticatedUser) (string, error) {
	if authenticatedUserHasActiveTenant(user) {
		return strings.TrimSpace(user.ActiveTenant.ID), nil
	}
	if user.IsPlatformAdmin {
		return "", nil
	}
	return "", errors.New("active tenant membership required")
}

func issueAuthSession(ctx context.Context, db *sql.DB, userID string, tenantID string, now time.Time, cfg authConfig, auditCtx requestAuditContext) (authIssuedTokens, error) {
	tenantID = strings.TrimSpace(tenantID)
	tokens, err := issueRawTokens(now, cfg)
	if err != nil {
		return authIssuedTokens{}, err
	}

	var sessionID string
	err = db.QueryRowContext(ctx, `
WITH new_session AS (
	INSERT INTO app_auth_sessions (user_id, tenant_id, expires_at, last_used_at, ip_address, user_agent)
	VALUES ($1::uuid, nullif($2, '')::uuid, $3, $4, nullif($5, '')::inet, $6)
	RETURNING id, user_id
), new_access AS (
	INSERT INTO app_auth_access_tokens (session_id, user_id, token_hash, issued_at, expires_at)
	SELECT id, user_id, $7, $4, $8 FROM new_session
), new_refresh AS (
	INSERT INTO app_auth_refresh_tokens (session_id, user_id, token_hash, issued_at, expires_at)
	SELECT id, user_id, $9, $4, $3 FROM new_session
)
SELECT id::text FROM new_session`,
		userID,
		tenantID,
		tokens.RefreshExpiresAt,
		now,
		auditCtx.IPAddress,
		auditCtx.UserAgent,
		tokens.AccessTokenHash,
		tokens.AccessExpiresAt,
		tokens.RefreshTokenHash,
	).Scan(&sessionID)
	if err != nil {
		return authIssuedTokens{}, err
	}
	return tokens, nil
}

func recordAuthLastLogin(ctx context.Context, db *sql.DB, userID string, now time.Time) error {
	_, err := db.ExecContext(ctx, `UPDATE app_users SET last_login_at = $2 WHERE id = $1::uuid`, userID, now)
	return err
}

func loadAuthSessionByAccessToken(ctx context.Context, db *sql.DB, accessTokenHash string, now time.Time) (authSessionResponse, error) {
	if accessTokenHash == "" {
		return authSessionResponse{}, sql.ErrNoRows
	}
	var session authSessionResponse
	var accessTokenID string
	var trialEndsAt sql.NullTime
	var currentPeriodEndsAt sql.NullTime
	err := db.QueryRowContext(ctx, `
SELECT at.id::text,
	at.expires_at,
	rt.expires_at,
	u.id::text,
	COALESCE(u.email, ''),
	u.phone,
	u.display_name,
	u.status,
	COALESCE(tenant.id::text, ''),
	tenant.code,
	tenant.name,
	COALESCE(tenant.status, ''),
	COALESCE(membership.status, ''),
	COALESCE(membership.is_owner, false),
	COALESCE(ts.status, ''),
	COALESCE(plan.code, ''),
	COALESCE(plan.name, ''),
	ts.trial_ends_at,
	ts.current_period_ends_at
FROM app_auth_access_tokens at
JOIN app_auth_sessions s ON s.id = at.session_id
JOIN app_auth_refresh_tokens rt ON rt.session_id = s.id AND rt.used_at IS NULL AND rt.revoked_at IS NULL
JOIN app_users u ON u.id = at.user_id
LEFT JOIN tenants tenant ON tenant.id = s.tenant_id
LEFT JOIN tenant_memberships membership ON membership.tenant_id = s.tenant_id
	AND membership.user_id = u.id
	AND membership.status = 'active'
LEFT JOIN tenant_subscriptions ts ON ts.tenant_id = tenant.id
LEFT JOIN subscription_plans plan ON plan.id = ts.plan_id
WHERE at.token_hash = $1
	AND at.revoked_at IS NULL
	AND at.expires_at > $2
	AND s.revoked_at IS NULL
	AND s.expires_at > $2
	AND rt.expires_at > $2
	AND u.status = 'active'
	AND (
		s.tenant_id IS NULL
		OR (
			tenant.status IN ('active', 'trial')
			AND membership.status = 'active'
		)
	)
ORDER BY rt.issued_at DESC
LIMIT 1`, accessTokenHash, now).Scan(
		&accessTokenID,
		&session.AccessExpiresAt,
		&session.RefreshExpiresAt,
		&session.User.ID,
		&session.User.Email,
		&session.User.Phone,
		&session.User.DisplayName,
		&session.User.Status,
		&session.User.ActiveTenant.ID,
		&session.User.ActiveTenant.Code,
		&session.User.ActiveTenant.Name,
		&session.User.ActiveTenant.Status,
		&session.User.ActiveTenant.MembershipStatus,
		&session.User.ActiveTenant.IsOwner,
		&session.User.ActiveTenant.SubscriptionStatus,
		&session.User.ActiveTenant.PlanCode,
		&session.User.ActiveTenant.PlanName,
		&trialEndsAt,
		&currentPeriodEndsAt,
	)
	if err != nil {
		return authSessionResponse{}, err
	}
	session.User.ActiveTenant.TrialEndsAt = formatNullDate(trialEndsAt)
	session.User.ActiveTenant.CurrentPeriodEndsAt = formatNullDate(currentPeriodEndsAt)
	if err := enrichAuthenticatedUser(ctx, db, &session.User); err != nil {
		return authSessionResponse{}, err
	}
	if session.User.ActiveTenant.ID == "" && !session.User.IsPlatformAdmin {
		return authSessionResponse{}, sql.ErrNoRows
	}
	_, _ = db.ExecContext(ctx, `
UPDATE app_auth_access_tokens
SET last_used_at = $2
WHERE id = $1::uuid`, accessTokenID, now)
	return session, nil
}

func rotateRefreshToken(ctx context.Context, db *sql.DB, refreshTokenHash string, now time.Time, cfg authConfig, auditCtx requestAuditContext) (authSessionResponse, authIssuedTokens, error) {
	tokens, err := issueRawTokens(now, cfg)
	if err != nil {
		return authSessionResponse{}, authIssuedTokens{}, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return authSessionResponse{}, authIssuedTokens{}, err
	}
	defer tx.Rollback()

	record, err := loadRefreshTokenForUpdate(ctx, tx, refreshTokenHash)
	if err != nil {
		return authSessionResponse{}, authIssuedTokens{}, err
	}
	if err := validateRefreshTokenRecord(record, now); err != nil {
		_, _ = tx.ExecContext(ctx, `
UPDATE app_auth_sessions
SET revoked_at = COALESCE(revoked_at, $2),
	revoked_reason = CASE WHEN revoked_reason = '' THEN $3 ELSE revoked_reason END
WHERE id = $1::uuid`, record.SessionID, now, "invalid refresh token reuse")
		_ = tx.Commit()
		return authSessionResponse{}, authIssuedTokens{}, err
	}

	var newRefreshTokenID string
	err = tx.QueryRowContext(ctx, `
WITH new_access AS (
	INSERT INTO app_auth_access_tokens (session_id, user_id, token_hash, issued_at, expires_at)
	VALUES ($1::uuid, $2::uuid, $3, $4, $5)
), new_refresh AS (
	INSERT INTO app_auth_refresh_tokens (session_id, user_id, token_hash, issued_at, expires_at)
	VALUES ($1::uuid, $2::uuid, $6, $4, $7)
	RETURNING id
)
SELECT id::text FROM new_refresh`,
		record.SessionID,
		record.UserID,
		tokens.AccessTokenHash,
		now,
		tokens.AccessExpiresAt,
		tokens.RefreshTokenHash,
		tokens.RefreshExpiresAt,
	).Scan(&newRefreshTokenID)
	if err != nil {
		return authSessionResponse{}, authIssuedTokens{}, err
	}
	_, err = tx.ExecContext(ctx, `
UPDATE app_auth_refresh_tokens
SET used_at = $2,
	replaced_by_token_id = $3::uuid
WHERE id = $1::uuid`, record.ID, now, newRefreshTokenID)
	if err != nil {
		return authSessionResponse{}, authIssuedTokens{}, err
	}
	_, err = tx.ExecContext(ctx, `
UPDATE app_auth_sessions
SET last_used_at = $2,
	ip_address = COALESCE(nullif($3, '')::inet, ip_address),
	user_agent = CASE WHEN $4 <> '' THEN $4 ELSE user_agent END
WHERE id = $1::uuid`, record.SessionID, now, auditCtx.IPAddress, auditCtx.UserAgent)
	if err != nil {
		return authSessionResponse{}, authIssuedTokens{}, err
	}

	session := authSessionResponse{
		User: authenticatedUser{
			ID:          record.UserID,
			Email:       record.UserEmail,
			Phone:       record.UserPhone,
			DisplayName: record.UserDisplayName,
			Status:      record.UserStatus,
			ActiveTenant: authTenantSummary{
				ID:                  record.TenantID,
				Code:                record.TenantCode,
				Name:                record.TenantName,
				Status:              record.TenantStatus,
				MembershipStatus:    record.MembershipStatus,
				IsOwner:             record.TenantIsOwner,
				SubscriptionStatus:  record.SubscriptionStatus,
				PlanCode:            record.PlanCode,
				PlanName:            record.PlanName,
				TrialEndsAt:         formatNullDate(record.TrialEndsAt),
				CurrentPeriodEndsAt: formatNullDate(record.CurrentPeriodEndsAt),
			},
		},
		AccessExpiresAt:  tokens.AccessExpiresAt,
		RefreshExpiresAt: tokens.RefreshExpiresAt,
	}
	if err := enrichAuthenticatedUser(ctx, db, &session.User); err != nil {
		return authSessionResponse{}, authIssuedTokens{}, err
	}
	schoolCount := -1
	if record.TenantID != "" {
		schoolCount, err = loadTenantSchoolCount(ctx, tx, record.TenantID)
		if err != nil {
			return authSessionResponse{}, authIssuedTokens{}, err
		}
	}
	session.Onboarding = buildAuthOnboardingSummary(session.User, schoolCount)
	if err := tx.Commit(); err != nil {
		return authSessionResponse{}, authIssuedTokens{}, err
	}
	return session, tokens, nil
}

func loadRefreshTokenForUpdate(ctx context.Context, tx *sql.Tx, refreshTokenHash string) (authTokenRecord, error) {
	var record authTokenRecord
	err := tx.QueryRowContext(ctx, `
SELECT rt.id::text,
	rt.session_id::text,
	rt.user_id::text,
	rt.expires_at,
	rt.used_at,
	rt.revoked_at,
	s.expires_at,
	s.revoked_at,
	COALESCE(u.email, ''),
	u.phone,
	u.display_name,
	u.status,
	COALESCE(s.tenant_id::text, ''),
	COALESCE(tenant.code, ''),
	COALESCE(tenant.name, ''),
	COALESCE(tenant.status, ''),
	COALESCE(membership.status, ''),
	COALESCE(membership.is_owner, false),
	COALESCE(ts.status, ''),
	COALESCE(plan.code, ''),
	COALESCE(plan.name, ''),
	ts.trial_ends_at,
	ts.current_period_ends_at,
	EXISTS (
		SELECT 1
		FROM app_user_roles ur
		JOIN app_roles r ON r.id = ur.role_id
		WHERE ur.user_id = u.id
			AND r.code = 'platform_admin'
	)
FROM app_auth_refresh_tokens rt
JOIN app_auth_sessions s ON s.id = rt.session_id
JOIN app_users u ON u.id = rt.user_id
LEFT JOIN tenants tenant ON tenant.id = s.tenant_id
LEFT JOIN tenant_memberships membership ON membership.tenant_id = s.tenant_id
	AND membership.user_id = u.id
LEFT JOIN tenant_subscriptions ts ON ts.tenant_id = tenant.id
LEFT JOIN subscription_plans plan ON plan.id = ts.plan_id
WHERE rt.token_hash = $1
FOR UPDATE OF rt`, refreshTokenHash).Scan(
		&record.ID,
		&record.SessionID,
		&record.UserID,
		&record.ExpiresAt,
		&record.UsedAt,
		&record.RevokedAt,
		&record.SessionExpiresAt,
		&record.SessionRevokedAt,
		&record.UserEmail,
		&record.UserPhone,
		&record.UserDisplayName,
		&record.UserStatus,
		&record.TenantID,
		&record.TenantCode,
		&record.TenantName,
		&record.TenantStatus,
		&record.MembershipStatus,
		&record.TenantIsOwner,
		&record.SubscriptionStatus,
		&record.PlanCode,
		&record.PlanName,
		&record.TrialEndsAt,
		&record.CurrentPeriodEndsAt,
		&record.IsPlatformAdmin,
	)
	return record, err
}

func revokeSessionByTokenHashes(ctx context.Context, db *sql.DB, accessTokenHash string, refreshTokenHash string, reason string) error {
	_, err := db.ExecContext(ctx, `
UPDATE app_auth_sessions
SET revoked_at = COALESCE(revoked_at, now()),
	revoked_reason = CASE WHEN revoked_reason = '' THEN $3 ELSE revoked_reason END
WHERE id IN (
	SELECT session_id FROM app_auth_access_tokens WHERE token_hash = $1
	UNION
	SELECT session_id FROM app_auth_refresh_tokens WHERE token_hash = $2
)`, accessTokenHash, refreshTokenHash, strings.TrimSpace(reason))
	return err
}
