package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const passwordResetTokenTTL = 2 * time.Hour

type passwordResetRequestInput struct {
	Identifier string `json:"identifier"`
	Email      string `json:"email"`
}

type passwordResetConfirmInput struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type passwordResetTokenResult struct {
	Token    string
	UserID   string
	Email    string
	TenantID string
}

func handleAuthPasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input passwordResetRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	identifier := normalizeAuthIdentifier(firstNonEmpty(input.Identifier, input.Email))
	if identifier == "" {
		http.Error(w, "email or phone is required", http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	response := map[string]any{
		"ok":      true,
		"message": "Nếu tài khoản tồn tại, hệ thống đã tạo hướng dẫn đặt lại mật khẩu.",
	}
	reset, err := createPasswordResetToken(r, db, identifier)
	if err != nil {
		http.Error(w, "cannot create password reset request", http.StatusInternalServerError)
		return
	}
	if reset.Token != "" && reset.Email != "" {
		_ = sendPasswordResetEmail(r, db, reset)
	}
	env, _ := loadAppEnvironment(os.LookupEnv)
	if reset.Token != "" && env != appEnvProduction {
		response["resetToken"] = reset.Token
		response["message"] = "Token đặt lại mật khẩu đã được tạo cho môi trường thử nghiệm."
	}
	writeJSON(w, http.StatusOK, response)
}

func handleAuthPasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input passwordResetConfirmInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input.Token = strings.TrimSpace(input.Token)
	input.Password = strings.TrimSpace(input.Password)
	if input.Token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}
	if err := validateAuthPassword(input.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	if err := confirmPasswordResetToken(r, db, input.Token, input.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func createPasswordResetToken(r *http.Request, db *sql.DB, identifier string) (passwordResetTokenResult, error) {
	var result passwordResetTokenResult
	err := db.QueryRowContext(r.Context(), `
SELECT u.id::text,
	COALESCE(u.email, ''),
	COALESCE((
		SELECT tm.tenant_id::text
		FROM tenant_memberships tm
		JOIN tenants t ON t.id = tm.tenant_id
		WHERE tm.user_id = u.id
			AND tm.status = 'active'
			AND t.status IN ('active', 'trial', 'suspended')
		ORDER BY tm.is_owner DESC, tm.created_at ASC
		LIMIT 1
	), '')
FROM app_users
u
WHERE status = 'active'
	AND (
		($1 <> '' AND lower(COALESCE(email, '')) = lower($1))
		OR ($2 <> '' AND phone = $2)
	)
LIMIT 1`, identifier, identifier).Scan(&result.UserID, &result.Email, &result.TenantID)
	if err == sql.ErrNoRows {
		return result, nil
	}
	if err != nil {
		return result, err
	}
	token, err := generateAuthToken()
	if err != nil {
		return result, err
	}
	_, err = db.ExecContext(r.Context(), `
INSERT INTO app_password_reset_tokens (user_id, token_hash, expires_at, user_agent)
VALUES ($1::uuid, $2, $3, $4)`,
		result.UserID,
		hashAuthToken(token),
		time.Now().UTC().Add(passwordResetTokenTTL),
		strings.TrimSpace(r.UserAgent()),
	)
	if err != nil {
		return result, err
	}
	result.Token = token
	return result, nil
}

func sendPasswordResetEmail(r *http.Request, db *sql.DB, reset passwordResetTokenResult) error {
	cfg, err := loadTenantEmailConfig(r.Context(), db, reset.TenantID)
	if err != nil {
		return err
	}
	resetURL := passwordResetURL(r, cfg, reset.Token)
	subject := "Đặt lại mật khẩu DEKISUGI Finance Hub"
	textBody := "Bạn vừa yêu cầu đặt lại mật khẩu DEKISUGI Finance Hub.\nMở link sau để đặt lại mật khẩu trong vòng 2 giờ:\n" + resetURL + "\nNếu bạn không yêu cầu thao tác này, vui lòng bỏ qua email."
	htmlBody := `<!doctype html><html><body style="font-family:Arial,Helvetica,sans-serif;line-height:1.5;color:#172033;">
<p>Bạn vừa yêu cầu đặt lại mật khẩu <strong>DEKISUGI Finance Hub</strong>.</p>
<p><a href="` + html.EscapeString(resetURL) + `" style="display:inline-block;padding:10px 14px;background:#14532d;color:#fff;text-decoration:none;border-radius:6px;">Đặt lại mật khẩu</a></p>
<p>Link có hiệu lực trong 2 giờ. Nếu bạn không yêu cầu thao tác này, vui lòng bỏ qua email.</p>
</body></html>`
	_, err = sendSimpleEmail(r.Context(), cfg, reset.Email, subject, htmlBody, textBody)
	return err
}

func passwordResetURL(r *http.Request, cfg emailConfig, token string) string {
	values := url.Values{}
	values.Set("resetToken", token)
	return appBaseURL(r, cfg) + "/?" + values.Encode()
}

func confirmPasswordResetToken(r *http.Request, db *sql.DB, token string, password string) error {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return err
	}
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var tokenID string
	var userID string
	err = tx.QueryRowContext(r.Context(), `
SELECT id::text, user_id::text
FROM app_password_reset_tokens
WHERE token_hash = $1
	AND used_at IS NULL
	AND expires_at > now()
FOR UPDATE`, hashAuthToken(token)).Scan(&tokenID, &userID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("token đặt lại mật khẩu không hợp lệ hoặc đã hết hạn")
	}
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(r.Context(), `
UPDATE app_users
SET password_hash = $2,
	password_updated_at = now(),
	updated_at = now()
WHERE id = $1::uuid`, userID, passwordHash); err != nil {
		return err
	}
	if _, err := tx.ExecContext(r.Context(), `
UPDATE app_password_reset_tokens
SET used_at = now()
WHERE id = $1::uuid`, tokenID); err != nil {
		return err
	}
	return tx.Commit()
}
