package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type platformPaymentProviderConfig struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	DisplayName       string `json:"displayName"`
	ProviderType      string `json:"providerType"`
	Status            string `json:"status"`
	Configured        bool   `json:"configured"`
	WebhookPath       string `json:"webhookPath,omitempty"`
	BankBIN           string `json:"bankBin,omitempty"`
	AccountNumber     string `json:"accountNumber,omitempty"`
	AccountName       string `json:"accountName,omitempty"`
	ClientID          string `json:"clientId,omitempty"`
	HasAPIKey         bool   `json:"hasApiKey"`
	APIKeyMasked      string `json:"apiKeyMasked,omitempty"`
	HasChecksumKey    bool   `json:"hasChecksumKey"`
	ChecksumKeyMasked string `json:"checksumKeyMasked,omitempty"`
	ReturnURL         string `json:"returnUrl,omitempty"`
	CancelURL         string `json:"cancelUrl,omitempty"`
	APIBaseURL        string `json:"apiBaseUrl,omitempty"`
	DefaultProvider   bool   `json:"defaultProvider"`
}

type platformPaymentProviderConfigSaveInput struct {
	TenantID      string `json:"tenantId"`
	Code          string `json:"code"`
	DisplayName   string `json:"displayName"`
	Status        string `json:"status"`
	BankBIN       string `json:"bankBin"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
	ClientID      string `json:"clientId"`
	APIKey        string `json:"apiKey"`
	ChecksumKey   string `json:"checksumKey"`
	ReturnURL     string `json:"returnUrl"`
	CancelURL     string `json:"cancelUrl"`
	APIBaseURL    string `json:"apiBaseUrl"`
	SetDefault    bool   `json:"setDefault"`
}

type platformTenantEmailConfigSaveInput struct {
	TenantID string `json:"tenantId"`
	emailConfig
}

type platformTenantEmailCronSaveInput struct {
	TenantID string `json:"tenantId"`
	emailCronRequest
}

func handlePlatformTenantPaymentProviders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlePlatformTenantPaymentProvidersGet(w, r)
	case http.MethodPost:
		handlePlatformTenantPaymentProviderSave(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePlatformTenantPaymentProvidersGet(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSpace(r.URL.Query().Get("tenantId"))
	if tenantID == "" {
		http.Error(w, "tenantId is required", http.StatusBadRequest)
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	if err := ensureDefaultPaymentProvidersForTenant(r.Context(), db, tenantID, ""); err != nil {
		http.Error(w, "cannot initialize tenant payment providers", http.StatusInternalServerError)
		return
	}
	providers, err := listPaymentProviders(r.Context(), db, tenantID)
	if err != nil {
		http.Error(w, "cannot load payment providers", http.StatusInternalServerError)
		return
	}
	defaultProvider, _ := loadTenantDefaultPaymentProviderCode(r.Context(), db, tenantID)
	writeJSON(w, http.StatusOK, map[string]any{"providers": platformPaymentProviderConfigs(providers, defaultProvider), "defaultProvider": defaultProvider})
}

func handlePlatformTenantPaymentProviderSave(w http.ResponseWriter, r *http.Request) {
	user, ok := authenticatedUserFromRequest(r)
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input platformPaymentProviderConfigSaveInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	input = normalizePlatformPaymentProviderConfigSaveInput(input)
	if err := validatePlatformPaymentProviderConfigSaveInput(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	if err := savePlatformPaymentProviderConfig(r.Context(), db, input, user.ID, auditContextFromRequest(r)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	providers, err := listPaymentProviders(r.Context(), db, input.TenantID)
	if err != nil {
		http.Error(w, "cannot reload payment providers", http.StatusInternalServerError)
		return
	}
	defaultProvider, _ := loadTenantDefaultPaymentProviderCode(r.Context(), db, input.TenantID)
	writeJSON(w, http.StatusOK, map[string]any{"providers": platformPaymentProviderConfigs(providers, defaultProvider), "defaultProvider": defaultProvider})
}

func handlePlatformTenantEmailConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tenantID := strings.TrimSpace(r.URL.Query().Get("tenantId"))
		if tenantID == "" {
			http.Error(w, "tenantId is required", http.StatusBadRequest)
			return
		}
		db, err := openMasterDataDatabase(r.Context())
		if err != nil {
			writeMasterDataDBError(w, err)
			return
		}
		defer db.Close()
		cfg, err := loadTenantEmailConfig(r.Context(), db, tenantID)
		if err != nil {
			http.Error(w, "cannot load tenant email config", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, emailConfigPublic(cfg))
	case http.MethodPost:
		user, ok := authenticatedUserFromRequest(r)
		if !ok {
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		var input platformTenantEmailConfigSaveInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		input.TenantID = strings.TrimSpace(input.TenantID)
		if input.TenantID == "" {
			http.Error(w, "tenantId is required", http.StatusBadRequest)
			return
		}
		db, err := openMasterDataDatabase(r.Context())
		if err != nil {
			writeMasterDataDBError(w, err)
			return
		}
		defer db.Close()
		current, _ := loadTenantEmailConfig(r.Context(), db, input.TenantID)
		cfg := normalizeEmailConfig(input.emailConfig)
		if cfg.APIKey == "" {
			cfg.APIKey = current.APIKey
		}
		if cfg.GmailAppPassword == "" {
			cfg.GmailAppPassword = current.GmailAppPassword
		}
		if err := saveTenantEmailConfig(r.Context(), db, input.TenantID, cfg, user.ID); err != nil {
			http.Error(w, "cannot save tenant email config", http.StatusInternalServerError)
			return
		}
		auditCtx := auditContextFromRequest(r)
		auditCtx.TenantID = input.TenantID
		_ = insertAuditLog(r.Context(), db, auditLogInput{
			Context:    auditCtx,
			Action:     "platform.tenant.email_config.update",
			EntityType: "tenant_email_config",
			EntityID:   input.TenantID,
			Metadata: map[string]any{
				"provider": cfg.Provider,
			},
		})
		writeJSON(w, http.StatusOK, emailConfigPublic(cfg))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePlatformTenantEmailCron(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tenantID := strings.TrimSpace(r.URL.Query().Get("tenantId"))
		if tenantID == "" {
			http.Error(w, "tenantId is required", http.StatusBadRequest)
			return
		}
		state, err := loadEmailCronStateForTenant(r.Context(), tenantID)
		if err != nil {
			http.Error(w, "cannot load tenant email cron", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, emailCronPublic(state, time.Now()))
	case http.MethodPost:
		user, ok := authenticatedUserFromRequest(r)
		if !ok {
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<20)
		var input platformTenantEmailCronSaveInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		input.TenantID = strings.TrimSpace(input.TenantID)
		if input.TenantID == "" {
			http.Error(w, "tenantId is required", http.StatusBadRequest)
			return
		}
		if input.Rows != nil && len(*input.Rows) > maxRows {
			http.Error(w, fmt.Sprintf("too many rows, max is %d", maxRows), http.StatusBadRequest)
			return
		}
		state, err := loadEmailCronStateForTenant(r.Context(), input.TenantID)
		if err != nil {
			http.Error(w, "cannot load tenant email cron", http.StatusInternalServerError)
			return
		}
		state = applyEmailCronRequest(state, input.emailCronRequest, time.Now())
		if state.Enabled {
			cfg, _ := loadEmailConfigForTenant(r.Context(), input.TenantID)
			if err := validateEmailConfigForSend(cfg); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		if err := saveEmailCronStateForTenantByUser(r.Context(), input.TenantID, state, user.ID); err != nil {
			http.Error(w, "cannot save tenant email cron", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, emailCronPublic(state, time.Now()))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func ensureDefaultPaymentProvidersForTenant(ctx context.Context, exec masterDataExecutor, tenantID string, userID string) error {
	if _, err := exec.ExecContext(ctx, `
INSERT INTO payment_providers (tenant_id, code, display_name, provider_type, status, config, created_by_user_id, updated_by_user_id)
VALUES
	($1::uuid, 'manual_vietqr', 'Manual VietQR', 'manual_vietqr', 'active', '{"reconciliation":"manual"}'::jsonb, nullif($2, '')::uuid, nullif($2, '')::uuid),
	($1::uuid, 'sepay', 'SePay Webhook', 'bank_webhook', 'sandbox', '{"signature":"optional"}'::jsonb, nullif($2, '')::uuid, nullif($2, '')::uuid),
	($1::uuid, 'payos', 'PayOS Payment Link', 'payment_link', 'sandbox', '{}'::jsonb, nullif($2, '')::uuid, nullif($2, '')::uuid)
ON CONFLICT (tenant_id, code) DO NOTHING`, tenantID, strings.TrimSpace(userID)); err != nil {
		return err
	}
	_, err := exec.ExecContext(ctx, `
INSERT INTO tenant_payment_settings (tenant_id, default_provider_code, updated_by_user_id)
VALUES ($1::uuid, 'manual_vietqr', nullif($2, '')::uuid)
ON CONFLICT (tenant_id) DO NOTHING`, tenantID, strings.TrimSpace(userID))
	return err
}

func normalizePlatformPaymentProviderConfigSaveInput(input platformPaymentProviderConfigSaveInput) platformPaymentProviderConfigSaveInput {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Code = headerKey(input.Code)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.Status = headerKey(input.Status)
	if input.Status == "" {
		input.Status = "active"
	}
	input.BankBIN = strings.TrimSpace(input.BankBIN)
	input.AccountNumber = strings.TrimSpace(input.AccountNumber)
	input.AccountName = strings.TrimSpace(input.AccountName)
	input.ClientID = strings.TrimSpace(input.ClientID)
	input.APIKey = strings.TrimSpace(input.APIKey)
	input.ChecksumKey = strings.TrimSpace(input.ChecksumKey)
	input.ReturnURL = strings.TrimSpace(input.ReturnURL)
	input.CancelURL = strings.TrimSpace(input.CancelURL)
	input.APIBaseURL = strings.TrimRight(strings.TrimSpace(input.APIBaseURL), "/")
	return input
}

func validatePlatformPaymentProviderConfigSaveInput(input platformPaymentProviderConfigSaveInput) error {
	if input.TenantID == "" {
		return fmt.Errorf("tenantId is required")
	}
	switch input.Code {
	case paymentProviderManualVietQR, paymentProviderSePay, paymentProviderPayOS:
	default:
		return fmt.Errorf("provider must be manual_vietqr, sepay, or payos")
	}
	switch input.Status {
	case "active", "inactive", "sandbox":
	default:
		return fmt.Errorf("provider status must be active, inactive, or sandbox")
	}
	return nil
}

func savePlatformPaymentProviderConfig(ctx context.Context, db *sql.DB, input platformPaymentProviderConfigSaveInput, userID string, auditCtx requestAuditContext) error {
	if err := ensureDefaultPaymentProvidersForTenant(ctx, db, input.TenantID, userID); err != nil {
		return err
	}
	provider, err := loadPaymentProviderForPlatformConfig(ctx, db, input.TenantID, input.Code)
	if err != nil {
		return err
	}
	config := copyConfigMap(provider.config)
	setConfigValue(config, "bankBin", input.BankBIN)
	setConfigValue(config, "bank_bin", input.BankBIN)
	setConfigValue(config, "accountNumber", input.AccountNumber)
	setConfigValue(config, "account_number", input.AccountNumber)
	setConfigValue(config, "accountName", input.AccountName)
	setConfigValue(config, "account_name", input.AccountName)
	setNestedConfigValue(config, "collection", "bankBin", input.BankBIN)
	setNestedConfigValue(config, "collection", "accountNumber", input.AccountNumber)
	setNestedConfigValue(config, "collection", "accountName", input.AccountName)
	setConfigValue(config, "clientId", input.ClientID)
	setNestedConfigValue(config, "credentials", "clientId", input.ClientID)
	setConfigValue(config, "apiKey", input.APIKey)
	setNestedConfigValue(config, "credentials", "apiKey", input.APIKey)
	setConfigValue(config, "checksumKey", input.ChecksumKey)
	setNestedConfigValue(config, "credentials", "checksumKey", input.ChecksumKey)
	setNestedConfigValue(config, "webhook", "checksumKey", input.ChecksumKey)
	setConfigValue(config, "returnUrl", input.ReturnURL)
	setNestedConfigValue(config, "checkout", "returnUrl", input.ReturnURL)
	setConfigValue(config, "cancelUrl", input.CancelURL)
	setNestedConfigValue(config, "checkout", "cancelUrl", input.CancelURL)
	setConfigValue(config, "apiBaseUrl", input.APIBaseURL)
	setNestedConfigValue(config, "checkout", "apiBaseUrl", input.APIBaseURL)

	configData, err := json.Marshal(config)
	if err != nil {
		return err
	}
	displayName := firstNonEmpty(input.DisplayName, provider.DisplayName)
	result, err := db.ExecContext(ctx, `
UPDATE payment_providers
SET display_name = $3,
	status = $4,
	config = $5::jsonb,
	updated_by_user_id = nullif($6, '')::uuid,
	updated_at = now()
WHERE tenant_id = $1::uuid
	AND code = $2`, input.TenantID, input.Code, displayName, input.Status, string(configData), userID)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("payment provider %q is not available", input.Code)
	}
	if input.SetDefault {
		if err := saveTenantDefaultPaymentProviderCode(ctx, db, input.TenantID, input.Code, userID); err != nil {
			return err
		}
	}
	auditCtx.TenantID = input.TenantID
	_ = insertAuditLog(ctx, db, auditLogInput{
		Context:    auditCtx,
		Action:     "platform.tenant.payment_provider.update",
		EntityType: "payment_provider",
		EntityID:   provider.ID,
		Metadata: map[string]any{
			"provider": input.Code,
			"status":   input.Status,
		},
	})
	return nil
}

func loadPaymentProviderForPlatformConfig(ctx context.Context, db *sql.DB, tenantID string, code string) (paymentProvider, error) {
	var provider paymentProvider
	var tenantCode string
	var configBytes []byte
	err := db.QueryRowContext(ctx, `
SELECT pp.id::text, pp.code, pp.display_name, pp.provider_type, pp.status, pp.tenant_id::text, t.code, pp.config
FROM payment_providers pp
JOIN tenants t ON t.id = pp.tenant_id
WHERE pp.tenant_id = $1::uuid
	AND pp.code = $2`, strings.TrimSpace(tenantID), headerKey(code)).Scan(
		&provider.ID,
		&provider.Code,
		&provider.DisplayName,
		&provider.ProviderType,
		&provider.Status,
		&provider.tenantID,
		&tenantCode,
		&configBytes,
	)
	if err == sql.ErrNoRows {
		return provider, fmt.Errorf("payment provider %q is not available", code)
	}
	if err != nil {
		return provider, err
	}
	provider.config = decodeMetadata(configBytes)
	provider.Configured = isPaymentProviderConfigured(provider)
	if provider.Code == paymentProviderSePay || provider.Code == paymentProviderPayOS {
		provider.WebhookPath = "/api/v1/payments/webhooks/" + strings.ToLower(tenantCode) + "/" + provider.Code
	}
	return provider, nil
}

func platformPaymentProviderConfigs(providers []paymentProvider, defaultProvider string) []platformPaymentProviderConfig {
	defaultProvider = headerKey(defaultProvider)
	if defaultProvider == "" {
		defaultProvider = paymentProviderManualVietQR
	}
	out := make([]platformPaymentProviderConfig, 0, len(providers))
	for _, provider := range providers {
		cfg := loadPayOSConfig(provider)
		apiKey := firstNonEmpty(providerConfigString(provider.config, "apiKey", "api_key", "xApiKey"), providerNestedConfigString(provider.config, "credentials", "apiKey", "api_key"))
		checksumKey := firstNonEmpty(providerConfigString(provider.config, "checksumKey", "checksum_key"), providerNestedConfigString(provider.config, "credentials", "checksumKey", "checksum_key"), providerNestedConfigString(provider.config, "webhook", "checksumKey", "checksum_key"))
		out = append(out, platformPaymentProviderConfig{
			ID:                provider.ID,
			Code:              provider.Code,
			DisplayName:       provider.DisplayName,
			ProviderType:      provider.ProviderType,
			Status:            provider.Status,
			Configured:        provider.Configured,
			WebhookPath:       provider.WebhookPath,
			BankBIN:           firstNonEmpty(providerConfigString(provider.config, "bankBin", "bank_bin"), providerNestedConfigString(provider.config, "collection", "bankBin", "bank_bin")),
			AccountNumber:     firstNonEmpty(providerConfigString(provider.config, "accountNumber", "account_number"), providerNestedConfigString(provider.config, "collection", "accountNumber", "account_number")),
			AccountName:       firstNonEmpty(providerConfigString(provider.config, "accountName", "account_name"), providerNestedConfigString(provider.config, "collection", "accountName", "account_name")),
			ClientID:          cfg.ClientID,
			HasAPIKey:         apiKey != "",
			APIKeyMasked:      maskAPIKey(apiKey),
			HasChecksumKey:    checksumKey != "",
			ChecksumKeyMasked: maskAPIKey(checksumKey),
			ReturnURL:         cfg.ReturnURL,
			CancelURL:         cfg.CancelURL,
			APIBaseURL:        cfg.APIBaseURL,
			DefaultProvider:   provider.Code == defaultProvider,
		})
	}
	return out
}

func loadTenantDefaultPaymentProviderCode(ctx context.Context, db *sql.DB, tenantID string) (string, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return paymentProviderManualVietQR, fmt.Errorf("tenant id is required")
	}
	var code string
	err := db.QueryRowContext(ctx, `
SELECT default_provider_code
FROM tenant_payment_settings
WHERE tenant_id = $1::uuid`, tenantID).Scan(&code)
	if err == sql.ErrNoRows {
		return paymentProviderManualVietQR, nil
	}
	if err != nil {
		return paymentProviderManualVietQR, err
	}
	code = headerKey(code)
	if code == "" {
		code = paymentProviderManualVietQR
	}
	return code, nil
}

func saveTenantDefaultPaymentProviderCode(ctx context.Context, db *sql.DB, tenantID string, providerCode string, userID string) error {
	providerCode = headerKey(providerCode)
	switch providerCode {
	case paymentProviderManualVietQR, paymentProviderSePay, paymentProviderPayOS:
	default:
		return fmt.Errorf("provider must be manual_vietqr, sepay, or payos")
	}
	_, err := db.ExecContext(ctx, `
INSERT INTO tenant_payment_settings (tenant_id, default_provider_code, updated_by_user_id)
VALUES ($1::uuid, $2, nullif($3, '')::uuid)
ON CONFLICT (tenant_id) DO UPDATE
SET default_provider_code = EXCLUDED.default_provider_code,
	updated_by_user_id = EXCLUDED.updated_by_user_id,
	updated_at = now()`, strings.TrimSpace(tenantID), providerCode, strings.TrimSpace(userID))
	return err
}

func copyConfigMap(input map[string]any) map[string]any {
	out := map[string]any{}
	for key, value := range input {
		nested, ok := value.(map[string]any)
		if ok {
			out[key] = copyConfigMap(nested)
			continue
		}
		out[key] = value
	}
	return out
}

func setConfigValue(config map[string]any, key string, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	config[key] = strings.TrimSpace(value)
}

func setNestedConfigValue(config map[string]any, parent string, key string, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	nested, _ := config[parent].(map[string]any)
	if nested == nil {
		nested = map[string]any{}
		config[parent] = nested
	}
	nested[key] = strings.TrimSpace(value)
}

func loadEmailConfigForTenant(ctx context.Context, tenantID string) (emailConfig, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return loadEmailConfig()
	}
	db, err := openMasterDataDatabase(ctx)
	if err != nil {
		return loadEmailConfig()
	}
	defer db.Close()
	return loadTenantEmailConfig(ctx, db, tenantID)
}

func loadTenantEmailConfig(ctx context.Context, db *sql.DB, tenantID string) (emailConfig, error) {
	var data []byte
	err := db.QueryRowContext(ctx, `
SELECT config
FROM tenant_email_configs
WHERE tenant_id = $1::uuid`, strings.TrimSpace(tenantID)).Scan(&data)
	if err == sql.ErrNoRows {
		return loadEmailConfig()
	}
	if err != nil {
		return emailConfig{}, err
	}
	cfg := defaultEmailConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return normalizeEmailConfig(cfg), nil
}

func saveTenantEmailConfig(ctx context.Context, db *sql.DB, tenantID string, cfg emailConfig, userID string) error {
	cfg = normalizeEmailConfig(cfg)
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, `
INSERT INTO tenant_email_configs (tenant_id, config, updated_by_user_id)
VALUES ($1::uuid, $2::jsonb, nullif($3, '')::uuid)
ON CONFLICT (tenant_id) DO UPDATE
SET config = EXCLUDED.config,
	updated_by_user_id = EXCLUDED.updated_by_user_id,
	updated_at = now()`, strings.TrimSpace(tenantID), string(data), strings.TrimSpace(userID))
	return err
}
